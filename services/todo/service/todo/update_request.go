package todo

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

type UpdateRequest struct {
	Title string `json:"title"`
	Description string `json:"description"`
	IsDone string `json:"is_done"`
	IsFavorite string `json:"is_favorite"`
}

func (c *UpdateRequest) Validate() error {
	validate := validator.New()
	if err := validate.Var(c.Title, "required,min=3,max=100"); err != nil {
		return errors.New("title harus diisi dan terdiri dari 3 s/d 100 karakter")
	}

	if err := validate.Var(c.Description, "required,min=10"); err != nil {
		return errors.New("description harus diisi dan minimal terdiri dari 10 karakter")
	}

	if err := validate.Var(c.IsDone, "omitempty,oneof=true false"); err != nil {
		return errors.New("is_done must between true or false")
	}

	if err := validate.Var(c.IsFavorite, "omitempty,oneof=true false"); err != nil {
		return errors.New("is_favorite must between true or false")
	}

	return nil
}

type DoneRequest struct {
	IsDone string `json:"is_done"`
}

func (c *DoneRequest) Validate() error {
	validate := validator.New()
	if err := validate.Var(c.IsDone, "oneof=true false"); err != nil {
		return errors.New("is_done must between true or false")
	}

	return nil
}

type FavoriteRequest struct {
	IsFavorite string `json:"is_favorite"`
}

func (c *FavoriteRequest) Validate() error {
	validate := validator.New()
	if err := validate.Var(c.IsFavorite, "oneof=true false"); err != nil {
		return errors.New("is_favorite must between true of false")
	}

	return nil
}

