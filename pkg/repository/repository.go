package service

type Authorization interface {

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

func NewService(db *sqlx.DB) *Repository {
	return &Repository{}
}