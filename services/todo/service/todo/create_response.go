package todo

import "github.com/jinzhu/gorm"

type CreateResponse struct {
	gorm.Model
	Title string `json:"title"`
	Description string `json:"description"`
	IsDone bool `json:"is_done"`
	IsFavorite bool `json:"is_favorite"`
}
