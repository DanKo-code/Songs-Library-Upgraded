package dtos

import (
	"SongsLibrary/internal/song/constants"
)

type GetSongsDTO struct {
	Id          string `form:"id" binding:"omitempty"`
	Name        string `form:"name" binding:"omitempty,max=100"`
	GroupName   string `form:"group_name" binding:"omitempty,max=100"`
	ReleaseDate string `form:"release_date" validate:"DateValidation"`
	Text        string `form:"text" binding:"omitempty,max=10000"`
	Link        string `form:"link" binding:"omitempty,url"`
	Page        int    `form:"page" binding:"omitempty,min=1"`
	PageSize    int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}

func (dto *GetSongsDTO) SetDefaults() {
	if dto.Page == 0 {
		dto.Page = constants.DefaultSongsPage
	}
	if dto.PageSize == 0 {
		dto.PageSize = constants.DefaultSongsPageSize
	}
}
