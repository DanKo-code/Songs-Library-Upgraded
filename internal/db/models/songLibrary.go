package models

import (
	"github.com/google/uuid"
	"time"
)

type Song struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	Name        string    `gorm:"size:255; unique"`
	GroupName   string    `gorm:"size:255"`
	ReleaseDate time.Time
	Text        string `gorm:"type:text"`
	Link        string `gorm:"type:text"`
}
