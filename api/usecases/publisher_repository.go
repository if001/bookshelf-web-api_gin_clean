package usecases

import "bookshelf-web-api_gin_clean/api/domain"

type PublisherRepository interface {
	FindAll() (*domain.CountedPublishers, error)
	Create(author domain.Publisher) (*domain.Publisher, error)
	Find(filter map[string]interface{}) (*domain.Publisher, error)
}
