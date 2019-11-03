package repository

import (
	"database/sql"
	"fmt"
	"github.com/ardiantirta/todo-crud/models"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type TodoRepository struct {
	Conn *sql.DB
}

func (t TodoRepository) Fetch(page int, limit int) (response []*models.Todo, count int, err error) {
	query := `Select id, title, description, completed, created_at, updated_at 
				from todo order by id desc limit $1 offset $2`

	queryCount := `Select count(id)
					from todo`

	start := time.Now()
	err = t.Conn.QueryRow(queryCount).Scan(&count)
	if err != nil {
		logrus.Error(err)
	}
	elapsed := time.Since(start)
	fmt.Printf("select count took %s\n", elapsed)

	start = time.Now()
	rows, err := t.Conn.Query(query, limit, (page-1)*limit)
	if err != nil {
		logrus.Error(err)
		return nil, 0, err
	}
	elapsed = time.Since(start)
	fmt.Printf("select query took %s\n", elapsed)

	defer func() {
		err = rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	start = time.Now()
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
	elapsed = time.Since(start)
	fmt.Printf("forloop insert to response took %s\n", elapsed)

	fmt.Println("count fetch", len(response))
	return response, count, nil
}

func (t TodoRepository) fetch(ch chan<- *models.Todo, wg *sync.WaitGroup, query string, page int, limit int, channel int) {
	defer wg.Done()
	rows, err := t.Conn.Query(query, limit/channel, (page-1)*limit)
	if err != nil {
		logrus.Error(err)
	}
	for rows.Next() {
		temp := new(models.Todo)
		_ = rows.Scan(
			&temp.ID,
			&temp.Title,
			&temp.Description,
			&temp.Completed,
			&temp.CreatedAt,
			&temp.UpdatedAt,
		)
		ch <- temp
		fmt.Println("write to channel")
	}

}

func (t TodoRepository) FetchWChannel(page int, limit int, channel int) (response []*models.Todo, count int, err error) {
	query := `Select id, title, description, completed, created_at, updated_at
				from todo order by id desc limit $1 offset $2`

	queryCount := `select count(id) from todo`

	start := time.Now()
	err = t.Conn.QueryRow(queryCount).Scan(&count)
	if err != nil {
		logrus.Error(err)
	}
	elapsed := time.Since(start)
	fmt.Printf("select count took %s\n", elapsed)

	ch := make(chan *models.Todo)
	var wg sync.WaitGroup

	start = time.Now()
	rows, err := t.Conn.Query(query, limit, (page-1)*limit)
	if err != nil {
		logrus.Error(err)
		return response, count, err
	}
	elapsed = time.Since(start)
	fmt.Printf("select rows took %s\n", elapsed)

	defer func() {
		err = rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	//for i := 1; i <= channel; i++ {
	wg.Add(1)
	go func() {
		fmt.Println("start goroutine")
		start = time.Now()
		defer wg.Done()
		for rows.Next() {
			temp := new(models.Todo)
			_ = rows.Scan(
				&temp.ID,
				&temp.Title,
				&temp.Description,
				&temp.Completed,
				&temp.CreatedAt,
				&temp.UpdatedAt,
			)
			ch <- temp
		}
		elapsed = time.Since(start)
		//fmt.Printf("scan each row took %s\n", elapsed)
	}()

	//}

	go func() {
		defer close(ch)
		wg.Wait()
		fmt.Println("stop goroutine")
	}()

	start = time.Now()
	for item := range ch {
		response = append(response, item)
	}
	elapsed = time.Since(start)
	fmt.Printf("forloop rows to response took %s\n", elapsed)

	fmt.Printf("count fetch %d\n", len(response))
	return response, count, err
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
	values ($1, $2, false, now(), now())`

	_, err := t.Conn.Exec(query, req.Title, req.Description)
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
