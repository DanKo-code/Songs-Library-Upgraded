package dtos

type GetSongsDTO struct {
	Name        string `form:"name" binding:"omitempty,max=100"`
	GroupName   string `form:"group_name" binding:"omitempty,max=100"`
	ReleaseDate string `form:"release_date" validate:"DateValidation"`
	Text        string `form:"text" binding:"omitempty,max=10000"`
	Link        string `form:"link" binding:"omitempty,url"`
	Page        int    `form:"page" binding:"omitempty,min=1"`
	PageSize    int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}
