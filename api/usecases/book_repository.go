package usecases

import (
	"bookshelf-web-api_gin_clean/api/domain"
)

type BookRepository interface {
	FindAll(filter map[string]interface{}, page uint64, perPage uint64, sortKey string) (*domain.PaginateBooks, error)
	Find(filter map[string]interface{}) (*domain.Book, error)
	Create(book domain.Book) (*domain.Book, error)
	Delete(filter map[string]interface{}) error
	Store(book domain.Book, filter map[string]interface{}) error
}
