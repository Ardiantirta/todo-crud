package repository

import (
	"context"
	"errors"
	"github.com/ardiantirta/todo-crud/services/todo/service/todo"
	"strconv"

	"github.com/jinzhu/gorm"
)

type Todo struct {
	gorm.Model
	Title string `json:"title" gorm:"type:varchar(255)"`
	Description string `json:"description" gorm:"type:text"`
	IsFavorite bool `json:"is_favorite"`
	IsDone bool `json:"is_done"`
}

type Repository interface {
	Create(ctx context.Context, data *Todo) (*todo.CreateResponse, error)
	Save(ctx context.Context, data *todo.ViewResponse) (*todo.ViewResponse, error)
	GetByID(ctx context.Context, todoID int) (*todo.ViewResponse, error)
	GetByTitle(ctx context.Context, params map[string]interface{}) ([]todo.ViewResponse, error)
	GetAll(ctx context.Context, params map[string]interface{}) ([]todo.ViewResponse, error)
	DeleteByID(ctx context.Context, todoID int) error
}

type TodoRepository struct {
	Conn *gorm.DB
}

func (c *TodoRepository) Create(ctx context.Context, data *Todo) (*todo.CreateResponse, error) {
	response := new(todo.CreateResponse)

	if err := c.Conn.Table("todos").Create(&data).Error; err != nil {
		return nil, errors.New("failed to create todo")
	}

	response.ID = data.ID
	response.Title = data.Title
	response.Description = data.Description
	response.IsDone = data.IsDone
	response.IsFavorite = data.IsFavorite
	response.CreatedAt = data.CreatedAt
	response.UpdatedAt = data.UpdatedAt
	response.DeletedAt = data.DeletedAt

	return response, nil
}

func (c *TodoRepository) Save(ctx context.Context, data *todo.ViewResponse) (*todo.ViewResponse, error) {
	if err := c.Conn.Table("todos").
		Save(&data).Error; err != nil {
			return nil, errors.New("failed to save todo")
	}

	return data, nil
}

func (c *TodoRepository) GetByID(ctx context.Context, todoID int) (*todo.ViewResponse, error) {
	response := new(todo.ViewResponse)

	data := new(Todo)
	if err := c.Conn.Table("todos").
		Where("id = ?", todoID).
		First(&data).Error; err != nil {
			return nil, errors.New("failed to get todo")
	}

	response.ID = data.ID
	response.Title = data.Title
	response.Description = data.Description
	response.IsDone = data.IsDone
	response.IsFavorite = data.IsFavorite
	response.CreatedAt = data.CreatedAt
	response.UpdatedAt = data.UpdatedAt
	response.DeletedAt = data.DeletedAt

	return response, nil
}

func (c *TodoRepository) GetByTitle(ctx context.Context, params map[string]interface{}) ([]todo.ViewResponse, error) {
	title := params["title"].(string)
	title = "%" + title + "%"

	response := make([]todo.ViewResponse, 0)
	if err := c.Conn.Table("todos").
		Where("title like ?", title).
		Find(&response).Error; err != nil {
			return nil, errors.New("failed to get todo")
	}

	return response, nil
}

func (c *TodoRepository) GetAll(ctx context.Context, params map[string]interface{}) ([]todo.ViewResponse, error) {
	response := make([]todo.ViewResponse, 0)

	isDoneStr := params["is_done"].(string)
	isFavoriteStr := params["is_favorite"].(string)

	db := c.Conn.Table("todos")

	if len(isDoneStr) >= 4  {
		isDoneBool, err := strconv.ParseBool(isDoneStr)
		if err != nil {

		} else {
			db  = db.Where("is_done = ?", isDoneBool)
		}
	}

	if len(isFavoriteStr) >= 4 {
		isFavoriteBool, err := strconv.ParseBool(params["is_favorite"].(string))
		if err != nil {

		} else {
			db = db.Where("is_favorite = ?", isFavoriteBool)
		}
	}

	if err := db.Find(&response).Error; err != nil {
		return nil, errors.New("failed to get todo")
	}

	return response, nil
}

func (c *TodoRepository) DeleteByID(ctx context.Context, todoID int) error {
	if err := c.Conn.Table("todos").
		Where("id = ?", todoID).
		Delete(Todo{}).Error; err != nil {
			return errors.New("failed to delete todo")
	}

	return nil
}

func NewTodoRepository(Conn *gorm.DB) Repository {
	return &TodoRepository{
		Conn: Conn,
	}
}
