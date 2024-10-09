package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type MusixMatchUseCase struct {
	baseURL       string
	getSongIPPath string
	getLyricsPath string
	apiKey        string
}

func CreateMusixMatchUseCase(baseURL, getSongIPPath, getLyricsPath, apiKey string) *MusixMatchUseCase {
	return &MusixMatchUseCase{
		baseURL:       baseURL,
		getSongIPPath: getSongIPPath,
		getLyricsPath: getLyricsPath,
		apiKey:        apiKey,
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

func (mmuc *MusixMatchUseCase) GetSongIP(groupName, song string) (string, string, string, error) {

	groupNameEscaped := url.QueryEscape(groupName)
	songEscaped := url.QueryEscape(song)

	musixMatchUrl := fmt.Sprintf(mmuc.baseURL+mmuc.getSongIPPath, groupNameEscaped, songEscaped, mmuc.apiKey)

	req, err := http.NewRequest("GET", musixMatchUrl, nil)
	if err != nil {
		return "", "", "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	var getSongIPResult GetSongIPResult
	if err := json.Unmarshal(body, &getSongIPResult); err != nil {
		return "", "", "", err
	}

	if len(getSongIPResult.Message.Body.TrackList) != 0 {
		songIp := getSongIPResult.Message.Body.TrackList[0].Track.Id
		link := getSongIPResult.Message.Body.TrackList[0].Track.Link
		releaseDate := getSongIPResult.Message.Body.TrackList[0].Track.ReleaseDate

		return strconv.Itoa(songIp), link, releaseDate, nil
	}

	return "", "", "", errors.New("the song was not found")
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

func (mmuc *MusixMatchUseCase) GetLyrics(ip string) (string, error) {

	musixMatchUrl := fmt.Sprintf(mmuc.baseURL+mmuc.getLyricsPath, ip, mmuc.apiKey)

	req, err := http.NewRequest("GET", musixMatchUrl, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	var getSongLyricsResult GetSongLyricsResult
	if err := json.Unmarshal(body, &getSongLyricsResult); err != nil {
		return "", err
	}

	lyricsBodyBorder := "\n...\n\n******* This Lyrics is NOT for Commercial use *******"
	lyricsBody := getSongLyricsResult.Message.Body.Lyrics.LyricsBody

	lyricsBody = strings.Split(lyricsBody, lyricsBodyBorder)[0]

	return lyricsBody, nil
}

/*func (mmuc *MusixMatchUseCase)
 */
