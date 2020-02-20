package todo

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

type CreateRequest struct {
	Title string `json:"title"`
	Description string `json:"description"`
}

func (c *CreateRequest) Validate() error {
	validate := validator.New()
	if err := validate.Var(c.Title, "required,min=3,max=100"); err != nil {
		return errors.New("title harus diisi dan terdiri dari 3 s/d 100 karakter")
	}

	if err := validate.Var(c.Description, "required,min=10"); err != nil {
		return errors.New("description harus diisi dan minimal terdiri dari 10 karakter")
	}

	return nil
}
