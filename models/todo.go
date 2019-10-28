package models

import "time"

// Todo represent the todo model
type Todo struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type ResponseHttp struct {
	Data    interface{} `json:"data"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
}
