package repository

import (
	"github.com/ardiantirta/todo-crud/models"
)

type Repository interface {
	Fetch(page int, limit int) ([]*models.Todo, int, error)
	GetById(id int64) (res *models.Todo, err error)
	Create(req models.Todo) error
	Update(id int64, req models.Todo) error
	Delete(id int64) error
}
