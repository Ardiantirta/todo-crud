package repository

import (
	"database/sql"
	"github.com/ardiantirta/todo-crud/models"
	"github.com/sirupsen/logrus"
)

type TodoRepository struct {
	Conn *sql.DB
}

func (t TodoRepository) Fetch(page int, limit int) (response []*models.Todo, count int, err error) {
	query := `Select id, title, description, completed, created_at, updated_at 
				from todo order by id limit $1 offset $2`

	queryCount := `Select count(id)
					from todo`

	err = t.Conn.QueryRow(queryCount).Scan(&count)
	if err != nil {
		logrus.Error(err)
	}

	rows, err := t.Conn.Query(query, limit, (page-1)*limit)
	if err != nil {
		logrus.Error(err)
		return nil, 0, err
	}

	defer func() {
		err = rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	for rows.Next() {
		temp := new(models.Todo)
		err = rows.Scan(
			&temp.ID,
			&temp.Title,
			&temp.Description,
			&temp.Completed,
			&temp.CreatedAt,
			&temp.UpdatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, 0, err
		}

		response = append(response, temp)
	}

	return response, count, nil
}

func (t TodoRepository) GetById(id int64) (*models.Todo, error) {
	query := `Select id, title, description, completed, created_at, updated_at
				from todo where id = $1`

	res := new(models.Todo)
	err := t.Conn.QueryRow(query, id).Scan(
		&res.ID,
		&res.Title,
		&res.Description,
		&res.Completed,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		logrus.Error(err)
	}

	return res, err
}

func (t TodoRepository) Create(req models.Todo) error {
	query := `insert into todo(title, description, completed, created_at, updated_at)
	values ($1, $2, $3, now(), now())`

	_, err := t.Conn.Exec(query, req.Title, req.Description, req.Completed)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (t TodoRepository) Update(id int64, req models.Todo) error {
	query := `select id, title, description, completed, created_at, updated_at
				from todo
				where id = $1`

	todo := new(models.Todo)
	err := t.Conn.QueryRow(query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)
	if err != nil {
		logrus.Error(err)
		return err
	}

	query = `update todo set 
                	title = $2,
                	description = $3,
                	completed = $4,
                	updated_at = now()
				where id = $1`

	if req.Title == "" {
		req.Title = todo.Title
	}

	if req.Description == "" {
		req.Description = todo.Description
	}

	_, err = t.Conn.Exec(query, id, req.Title, req.Description, req.Completed)
	if err != nil {
		logrus.Error(err)
		return err
	}

	return nil
}

func (t TodoRepository) Delete(id int64) error {
	query := `delete from todo
				where id = $1`

	_, err := t.Conn.Exec(query, id)
	if err != nil {
		logrus.Error(err)

		return err
	}

	return nil
}

func NewTodoRepository(Conn *sql.DB) Repository {
	return &TodoRepository{
		Conn: Conn,
	}
}
