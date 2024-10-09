package dtos

import "time"

type GetSongsDTO struct {
	Name        string    `form:"name"`
	GroupName   string    `form:"group_name"`
	ReleaseDate time.Time `form:"release_date"`
	Text        string    `form:"text"`
	Link        string    `form:"link"`
	Page        int       `form:"page"`
	PageSize    int       `form:"page_size"`
}

type CreateSongDTO struct {
	Group string `form:"group"`
	Song  string `form:"song"`
}
