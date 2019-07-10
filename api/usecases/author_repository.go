package usecases

import "bookshelf-web-api_gin_clean/api/domain"

type AuthorRepository interface {
	FindAll() (*domain.CountedAuthors, error)
	Create(author domain.Author) (*domain.Author, error)
}
