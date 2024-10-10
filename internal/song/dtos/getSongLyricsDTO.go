package dtos

import (
	"SongsLibrary/internal/song/constants"
	"github.com/google/uuid"
)

type GetSongLyricsDTO struct {
	Id       uuid.UUID `form:"id"`
	Page     int       `form:"page" binding:"omitempty,min=1"`
	PageSize int       `form:"page_size" binding:"omitempty,min=1,max=100"`
}

func (dto *GetSongLyricsDTO) SetDefaults() {
	if dto.Page == 0 {
		dto.Page = constants.DefaultLyricsPage
	}
	if dto.PageSize == 0 {
		dto.PageSize = constants.DefaultLyricsPageSize
	}
}
