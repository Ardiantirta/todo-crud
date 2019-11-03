package service

import (
	"github.com/ardiantirta/todo-crud/models"
	"github.com/ardiantirta/todo-crud/todo/repository"
	"github.com/sirupsen/logrus"
)

type TodoService struct {
	todoRepo repository.Repository
}

func (t *TodoService) Fetch(page int, limit int) (response []*models.Todo, count int, err error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 100
	}

	listTodo, count, err := t.todoRepo.Fetch(page, limit)
	if err != nil {
		logrus.Error(err)
		return nil, count, err
	}

	return listTodo, count, nil
}

func (t *TodoService) FetchWChannel(page int, limit int, channel int) (response []*models.Todo, count int, err error) {
	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 100
	}

	if channel == 0 {
		channel = 1
	}

	listTodo, count, err := t.todoRepo.FetchWChannel(page, limit, channel)
	if err != nil {
		logrus.Error(err)
		return nil, count, err
	}

	return listTodo, count, err
}

func (t *TodoService) GetById(id int64) (res *models.Todo, err error) {
	res, err = t.todoRepo.GetById(id)
	if err != nil {
		logrus.Error(err)
	}

	return res, err
}

func (t *TodoService) Create(req models.Todo) error {
	err := t.todoRepo.Create(req)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (t *TodoService) CreateBulk(req models.Todo, bulkCount int) error {
	err := t.todoRepo.CreateBulk(req, bulkCount)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (t *TodoService) Update(id int64, req models.Todo) error {
	err := t.todoRepo.Update(id, req)
	if err != nil {
		return err
	}
	return nil
}

func (t *TodoService) Delete(id int64) error {
	err := t.todoRepo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func NewTodoService(tr repository.Repository) Service {
	return &TodoService{
		todoRepo: tr,
	}
}
