package usecase

import (
	"SongsLibrary/internal/song"
	logrusCustom "SongsLibrary/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type MusixMatchUseCase struct {
	baseURL       string
	getSongIPPath string
	getLyricsPath string
	apiKey        string

	geniusBaseURL               string
	geniusGetSongReleaseDateURL string
	geniusAuthorization         string
}

func CreateMusixMatchUseCase(baseURL, getSongIPPath, getLyricsPath, apiKey, geniusBaseURL, geniusGetSongReleaseDateURL, geniusAuthorization string) *MusixMatchUseCase {
	return &MusixMatchUseCase{
		baseURL:       baseURL,
		getSongIPPath: getSongIPPath,
		getLyricsPath: getLyricsPath,
		apiKey:        apiKey,

		geniusBaseURL:               geniusBaseURL,
		geniusGetSongReleaseDateURL: geniusGetSongReleaseDateURL,
		geniusAuthorization:         geniusAuthorization,
	}
}

type GetSongIPResult struct {
	Message struct {
		Body struct {
			TrackList []TrackWrapper `json:"track_list"`
		} `json:"body"`
	} `json:"message"`
}

type TrackWrapper struct {
	Track Track `json:"track"` // Здесь содержится сам трек
}

type Track struct {
	Id          int    `json:"commontrack_id"`
	ReleaseDate string `json:"updated_time"`
	Link        string `json:"track_share_url"`
}

// Genius
type GetSongReleaseDateResult struct {
	Response struct {
		HitsList []HitWrapper `json:"hits"`
	} `json:"response"`
}

type HitWrapper struct {
	Result Result `json:"result"`
}

type Result struct {
	ReleaseDateComponents struct {
		Year  int `json:"year"`
		Month int `json:"month"`
		Day   int `json:"day"`
	} `json:"release_date_components"`
}

func (mmuc *MusixMatchUseCase) GetSongIP(ctx context.Context, groupName, songName string) (string, string, string, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered GetSongIP UseCase with parameters: groupName:%s, song:%s", groupName, songName))

	groupNameEscaped := url.QueryEscape(groupName)
	songEscaped := url.QueryEscape(songName)

	musixMatchUrl := fmt.Sprintf(mmuc.baseURL+mmuc.getSongIPPath, groupNameEscaped, songEscaped, mmuc.apiKey)
	geniusUrl := fmt.Sprintf(mmuc.geniusBaseURL+mmuc.geniusGetSongReleaseDateURL, songEscaped, groupNameEscaped)

	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Builded GetSongIP: musixMatchUrl:%s", musixMatchUrl))
	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Builded GetSongReleaseDate URL: geniusUrl:%s", geniusUrl))

	req, err := http.NewRequestWithContext(ctx, "GET", musixMatchUrl, nil)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return "", "", "", song.ErrorGetSongData
	}
	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Builded GetSongIP REQ: %+v", req))

	reqGenius, err := http.NewRequestWithContext(ctx, "GET", geniusUrl, nil)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return "", "", "", song.ErrorGetSongData
	}
	reqGenius.Header.Set("Authorization", mmuc.geniusAuthorization)
	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Builded GetSongReleaseDate REQ: %+v", req))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return "", "", "", song.ErrorGetSongData
	}
	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Recieved GetSongIP RES: %+v", resp))
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return "", "", "", song.ErrorGetSongData
	}

	if resp.StatusCode != http.StatusOK {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, fmt.Sprintf("unexpected status code: %d, body: %s", resp.StatusCode, body))

		return "", "", "", song.ErrorGetSongData
	}

	var getSongIPResult GetSongIPResult
	if err := json.Unmarshal(body, &getSongIPResult); err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return "", "", "", song.ErrorGetSongData
	}

	//genius
	clientGenius := &http.Client{}
	respGenius, err := clientGenius.Do(reqGenius)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return "", "", "", song.ErrorGetSongData
	}
	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Recieved GetSongReleaseDate RES: %+v", resp))
	defer respGenius.Body.Close()

	bodyGenius, err := io.ReadAll(respGenius.Body)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return "", "", "", song.ErrorGetSongData
	}

	if respGenius.StatusCode != http.StatusOK {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, fmt.Sprintf("unexpected status code: %d, body: %s", resp.StatusCode, body))

		return "", "", "", song.ErrorGetSongData
	}

	var getSongReleaseDateResult GetSongReleaseDateResult
	if err := json.Unmarshal(bodyGenius, &getSongReleaseDateResult); err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return "", "", "", song.ErrorGetSongData
	}

	if len(getSongIPResult.Message.Body.TrackList) != 0 {
		songIp := getSongIPResult.Message.Body.TrackList[0].Track.Id
		link := getSongIPResult.Message.Body.TrackList[0].Track.Link

		var releaseDate time.Time
		if respGenius.StatusCode == http.StatusOK {
			rdc := getSongReleaseDateResult.Response.HitsList[0].Result.ReleaseDateComponents

			releaseDate = time.Date(rdc.Year, time.Month(rdc.Month), rdc.Day, 0, 0, 0, 0, time.UTC)
		}

		logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Exiting GetSongIP UseCase with song data: ip:%s, link:%s, releaseDate:%s", strconv.Itoa(songIp), link, releaseDate))
		return strconv.Itoa(songIp), link, releaseDate.String(), nil
	}

	return "", "", "", song.ErrorGetSongData
}

type GetSongLyricsResult struct {
	Message struct {
		Body struct {
			Lyrics struct {
				LyricsBody string `json:"lyrics_body"`
			} `json:"lyrics"`
		} `json:"body"`
	} `json:"message"`
}

func (mmuc *MusixMatchUseCase) GetLyrics(ctx context.Context, ip string) (string, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered GetLyrics UseCase with parameter: ip:%s", ip))

	musixMatchUrl := fmt.Sprintf(mmuc.baseURL+mmuc.getLyricsPath, ip, mmuc.apiKey)
	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Builded GetSongLyrics URL: musixMatchUrl:%s", musixMatchUrl))

	req, err := http.NewRequestWithContext(ctx, "GET", musixMatchUrl, nil)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return "", song.ErrorGetSongLyrics
	}
	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Builded GetSongLyrics REQ: %+v", req))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return "", song.ErrorGetSongData
	}
	logrusCustom.LogWithLocation(logrus.DebugLevel, fmt.Sprintf("Recieved GetSongIP RES: %+v", resp))
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return "", song.ErrorGetSongLyrics
	}

	if resp.StatusCode != http.StatusOK {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, fmt.Sprintf("unexpected status code: %d, body: %s", resp.StatusCode, body))

		return "", song.ErrorGetSongLyrics
	}

	var getSongLyricsResult GetSongLyricsResult
	if err := json.Unmarshal(body, &getSongLyricsResult); err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, fmt.Sprintf("unexpected status code: %d, body: %s", resp.StatusCode, body))

		return "", song.ErrorGetSongLyrics
	}

	lyricsBodyBorder := "\n...\n\n******* This Lyrics is NOT for Commercial use *******"
	lyricsBody := getSongLyricsResult.Message.Body.Lyrics.LyricsBody

	lyricsBody = strings.Split(lyricsBody, lyricsBodyBorder)[0]

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Exiting GetLyrics UseCase with song lyrics: %+v", lyricsBody))

	return lyricsBody, nil
}
