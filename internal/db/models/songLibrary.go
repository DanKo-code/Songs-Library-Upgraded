package models

import (
	"github.com/google/uuid"
	"time"
)

type Author struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	GroupName string    `gorm:"size:255;unique"`
	Songs     []Song    `gorm:"foreignKey:AuthorID"`
}

type Song struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	Name        string    `gorm:"size:255;uniqueIndex:idx_song_author"`
	AuthorID    uuid.UUID `gorm:"not null;uniqueIndex:idx_song_author"`
	Author      Author    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ReleaseDate time.Time
	Text        string `gorm:"type:text"`
	Link        string `gorm:"type:text"`
}
