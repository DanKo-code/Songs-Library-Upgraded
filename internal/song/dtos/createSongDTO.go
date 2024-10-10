package dtos

type CreateSongDTO struct {
	Group string `form:"group"`
	Song  string `form:"song"`
}
