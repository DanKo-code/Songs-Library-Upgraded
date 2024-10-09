package postgres

import (
	"SongsLibrary/internal/db/models"
	"SongsLibrary/internal/song/dtos"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type SongRepository struct {
	db *gorm.DB
}

func NewSongRepository(db *gorm.DB) *SongRepository {
	return &SongRepository{db: db}
}

func (sr *SongRepository) GetSongs(gsdto *dtos.GetSongsDTO) ([]models.Song, error) {
	var songs []models.Song

	query := sr.db.Model(&models.Song{})

	if gsdto.Name != "" {
		query = query.Where("name LIKE ?", "%"+gsdto.Name+"%")
	}
	if gsdto.GroupName != "" {
		query = query.Where("group_name LIKE ?", "%"+gsdto.GroupName+"%")
	}
	if !gsdto.ReleaseDate.IsZero() {
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

	if err := query.Find(&songs).Error; err != nil {
		return nil, err
	}

	return songs, nil
}

func (sr *SongRepository) DeleteSong(id uuid.UUID) (*models.Song, error) {
	var songToDelete models.Song

	if err := sr.db.First(&songToDelete, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	if err := sr.db.Delete(&models.Song{}, id).Error; err != nil {
		return nil, err
	}

	return &songToDelete, nil
}

func (sr *SongRepository) UpdateSong(fieldsToUpdate *models.Song) (*models.Song, error) {
	if err := sr.db.Model(&models.Song{}).Where("id = ?", fieldsToUpdate.ID).Updates(fieldsToUpdate).Error; err != nil {
		return nil, err
	}

	return fieldsToUpdate, nil
}

func (sr *SongRepository) CreateSong(group, song, lyrics, link, releaseDate string) (*models.Song, error) {
	releaseDateCasted, err := time.Parse(time.RFC3339, releaseDate)
	if err != nil {
		return nil, err
	}

	var songToCreate *models.Song = &models.Song{ID: uuid.New(), Name: song, GroupName: group, Text: lyrics, Link: link, ReleaseDate: releaseDateCasted}
	if err := sr.db.Create(&songToCreate).Error; err != nil {
		return nil, err
	}

	return songToCreate, nil
}
