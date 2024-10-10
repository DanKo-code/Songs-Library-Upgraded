package dtos

import "github.com/google/uuid"

type GetSongLyricsDTO struct {
	Id       uuid.UUID `form:"id"`
	Page     int       `form:"page"`
	PageSize int       `form:"page_size"`
}
