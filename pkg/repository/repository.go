package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/DavelPurov777/todo-app-golang"
)

type Authorization interface {
	CreateUser(todo.User) (int, error)
}

type TodoList interface {

}

type TodoItem interface {

}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}
// поскольку репозиторий должен работать с базой данных, передадим db *sqlx.DB в качестве параметра функции
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
	}
}