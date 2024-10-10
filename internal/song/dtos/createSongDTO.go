package dtos

type CreateSongDTO struct {
	Group string `form:"group" binding:"omitempty,max=100"`
	Song  string `form:"song" binding:"omitempty,max=100"`
}
