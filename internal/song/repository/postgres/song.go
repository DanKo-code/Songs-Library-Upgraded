package postgres

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song"
	"SongsLibrary/internal/song/constants"
	"SongsLibrary/internal/song/dtos"
	logrusCustom "SongsLibrary/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

type SongRepository struct {
	db *gorm.DB
}

func NewSongRepository(db *gorm.DB) *SongRepository {
	return &SongRepository{db: db}
}

func (sr *SongRepository) GetSongs(ctx context.Context, gsdto *dtos.GetSongsDTO) ([]models.Song, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered GetSongs Repository with parameters: %+v", gsdto))

	var songs []models.Song

	query := sr.db.WithContext(ctx).Model(&models.Song{})

	if gsdto.Name != "" {
		query = query.Where("name LIKE ?", "%"+gsdto.Name+"%")
	}
	if gsdto.GroupName != "" {
		query = query.Where("group_name LIKE ?", "%"+gsdto.GroupName+"%")
	}
	if gsdto.ReleaseDate != "" {
		query = query.Where("release_date = ?", gsdto.ReleaseDate)
	}
	if gsdto.Text != "" {
		query = query.Where("text LIKE ?", "%"+gsdto.Text+"%")
	}
	if gsdto.Link != "" {
		query = query.Where("link LIKE ?", "%"+gsdto.Link+"%")
	}

	offset := (gsdto.Page - 1) * gsdto.PageSize
	query = query.Offset(offset).Limit(gsdto.PageSize)

	query = query.Debug()

	if err := query.Find(&songs).Error; err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		return nil, err
	}

	if len(songs) == 0 {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, song.SongsNotFound.Error())

		return nil, song.SongsNotFound
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Exiting GetSongs Repository with songs: %+v", songs))

	return songs, nil
}

func (sr *SongRepository) DeleteSong(ctx context.Context, id uuid.UUID) (*models.Song, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered DeleteSong Repository with parameter: %s", id.String()))

	var songToDelete models.Song

	if err := sr.db.WithContext(ctx).Debug().First(&songToDelete, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrusCustom.LogWithLocation(logrus.ErrorLevel, song.SongsNotFound.Error())

			return nil, song.SongsNotFound
		}

		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())
		return nil, err
	}

	if err := sr.db.WithContext(ctx).Debug().Delete(&models.Song{}, id).Error; err != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, song.SongsNotFound.Error())

		return nil, err
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Exiting DeleteSongs Repository with deleted song: %+v", songToDelete))

	return &songToDelete, nil
}

func (sr *SongRepository) UpdateSong(ctx context.Context, fieldsToUpdate *models.Song) (*models.Song, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered UpdateSong Repository with parameter: %+v", fieldsToUpdate))

	result := sr.db.WithContext(ctx).Debug().Model(&models.Song{}).Where("id = ?", fieldsToUpdate.ID).Updates(fieldsToUpdate)
	if result.Error != nil {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, result.Error.Error())

		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		logrusCustom.LogWithLocation(logrus.ErrorLevel, song.SongsNotFound.Error())

		return nil, song.SongsNotFound
	}

	var updatedSong models.Song
	if err := sr.db.WithContext(ctx).Debug().First(&updatedSong, "id = ?", fieldsToUpdate.ID).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			logrusCustom.LogWithLocation(logrus.ErrorLevel, song.SongsNotFound.Error())

			return nil, song.SongsNotFound
		}

		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())
		return nil, err
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Exiting UpdateSongs Repository with updated song: %+v", updatedSong))

	return &updatedSong, nil
}

func (sr *SongRepository) CreateSong(ctx context.Context, releaseDate time.Time, group string, songName string, lyrics string, link string) (*models.Song, error) {
	var songToCreate *models.Song = &models.Song{ID: uuid.New(), Name: songName, GroupName: group, Text: lyrics, Link: link, ReleaseDate: releaseDate}
	if err := sr.db.WithContext(ctx).Debug().Create(&songToCreate).Error; err != nil {

		logrusCustom.LogWithLocation(logrus.ErrorLevel, err.Error())

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == constants.DbUniqueConstrintErr {
			return nil, song.SongAlreadyExists
		}

		return nil, err
	}

	return songToCreate, nil
}

func (sr *SongRepository) GetSong(ctx context.Context, id uuid.UUID) (*models.Song, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered GetSong Repository with parameter: id:%s", id.String()))

	var songToGet models.Song
	if err := sr.db.WithContext(ctx).Debug().First(&songToGet, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, song.SongsNotFound
		}

		return nil, err
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Exiting GetSong Repository with songs: %+v", songToGet))

	return &songToGet, nil
}

func (sr *SongRepository) GetSongByName(ctx context.Context, name string) (*models.Song, error) {

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Entered GetSongByName Repository with parameter: name:%s", name))

	var songToGet models.Song
	if err := sr.db.WithContext(ctx).Debug().First(&songToGet, "name = ?", name).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, song.SongsNotFound
		}

		return nil, err
	}

	logrusCustom.LogWithLocation(logrus.InfoLevel, fmt.Sprintf("Exiting GetSongByName Repository with song: %+v", songToGet))

	return &songToGet, nil
}
