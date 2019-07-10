package usecases

import "bookshelf-web-api_gin_clean/api/domain"

type DescriptionRepository interface {
	FindAll(filter map[string]interface{}, page uint64, perPage uint64) (*domain.Descriptions, error)
	Find(filter map[string]interface{})  (*domain.Description, error)
	Create(description domain.Description) (desc *domain.Description, err error)
	Delete(description domain.Description) error
}
