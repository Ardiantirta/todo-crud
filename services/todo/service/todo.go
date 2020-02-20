package service

import (
	"context"
	"github.com/ardiantirta/todo-crud/services/todo/repository"
	"github.com/ardiantirta/todo-crud/services/todo/service/todo"
	"strconv"
)

type Service interface {
	Create(ctx context.Context, form *todo.CreateRequest) (*todo.CreateResponse, error)
	GetByID(ctx context.Context, id int) (*todo.ViewResponse, error)
	GetByTitle(ctx context.Context, params map[string]interface{}) ([]todo.ViewResponse, error)
	GetAll(ctx context.Context, params map[string]interface{}) ([]todo.ViewResponse, error)
	UpdateData(ctx context.Context, id int, form *todo.UpdateRequest) (*todo.ViewResponse, error)
	MarkAsDone(ctx context.Context, id int, form *todo.DoneRequest) error
	MarkAsFavorite(ctx context.Context,id int, form *todo.FavoriteRequest) error
	DeleteByID(ctx context.Context, id int) error
}

type TodoService struct {
	TodoRepository repository.Repository
}

func (c *TodoService) Create(ctx context.Context, form *todo.CreateRequest) (*todo.CreateResponse, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	data := &repository.Todo{
		Title:       form.Title,
		Description: form.Description,
		IsFavorite:  false,
		IsDone:      false,
	}

	response, err := c.TodoRepository.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *TodoService) GetByID(ctx context.Context, id int) (*todo.ViewResponse, error) {
	response, err := c.TodoRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *TodoService) GetByTitle(ctx context.Context, params map[string]interface{}) ([]todo.ViewResponse, error) {
	response, err := c.TodoRepository.GetByTitle(ctx, params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *TodoService) GetAll(ctx context.Context, params map[string]interface{}) ([]todo.ViewResponse, error) {
	response, err := c.TodoRepository.GetAll(ctx, params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *TodoService) UpdateData(ctx context.Context, id int, form *todo.UpdateRequest) (*todo.ViewResponse, error) {
	if err := form.Validate(); err != nil {
		return nil, err
	}

	isDone, _ := strconv.ParseBool(form.IsDone)
	isFavorite, _ := strconv.ParseBool(form.IsFavorite)

	response, err := c.TodoRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	response.Title = form.Title
	response.Description = form.Description
	response.IsDone = isDone
	response.IsFavorite = isFavorite

	response, err = c.TodoRepository.Save(ctx, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *TodoService) MarkAsDone(ctx context.Context, id int, form *todo.DoneRequest) error {
	if err := form.Validate(); err != nil {
		return err
	}

	isDone, _ := strconv.ParseBool(form.IsDone)

	response, err := c.TodoRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	response.IsDone = isDone

	_, err = c.TodoRepository.Save(ctx, response)
	if err != nil {
		return err
	}

	return nil
}

func (c *TodoService) MarkAsFavorite(ctx context.Context, id int, form *todo.FavoriteRequest) error {
	if err := form.Validate(); err != nil {
		return err
	}

	isFavorite, _ := strconv.ParseBool(form.IsFavorite)

	response, err := c.TodoRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	response.IsFavorite = isFavorite

	_, err = c.TodoRepository.Save(ctx, response)
	if err != nil {
		return err
	}

	return nil
}

func (c *TodoService) DeleteByID(ctx context.Context, id int) error {
	_, err := c.TodoRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := c.TodoRepository.DeleteByID(ctx, id); err != nil {
		return err
	}

	return nil
}

func NewTodoService(todoRepository repository.Repository) Service {
	return &TodoService{
		TodoRepository: todoRepository,
	}
}
