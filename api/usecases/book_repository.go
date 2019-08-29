package usecases

import (
	"bookshelf-web-api_gin_clean/api/domain"
)

type BookRepository interface {
	FindAll(filter map[string]interface{}, page uint64, perPage uint64, sortKey string) (*domain.PaginateBooks, error)
	Find(filter map[string]interface{}) (*domain.Book, error)
	Create(book domain.Book) (*domain.Book, error)
	Delete(filter map[string]interface{}) error
	Store(book domain.Book) error
	UpdateUpdatedAt(filter map[string]interface{}) error
	CountByAuthor(filter map[string]interface{}) (*domain.CountedNames, error)
	CountByPublisher(filter map[string]interface{}) (*domain.CountedNames, error)
	CountByDate(filter map[string]interface{}, key, format string) (*domain.CountedDates, error)
}
