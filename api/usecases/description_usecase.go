package usecases

import "bookshelf-web-api_gin_clean/api/domain"

type descriptionUseCase struct {
	DescriptionRepo DescriptionRepository
	BookRepository  BookRepository
}
type DescriptionUseCase interface {
	GetAllDescriptions(filter map[string]interface{}, page, perPage uint64) (*domain.Descriptions, error)
	GetDescription(filter map[string]interface{}) (*domain.Description, error)
	CreateDescription(createDescription domain.Description) (*domain.Description, error)
	UpdateDescription(updateDescription domain.Description, filter map[string]interface{}) (error)
	DeleteDescription(deleteDescription domain.Description) (error)
}

func NewDescriptionUseCase(descRepo DescriptionRepository, bookRepo BookRepository) DescriptionUseCase {
	return &descriptionUseCase{DescriptionRepo: descRepo, BookRepository: bookRepo}
}

func (b *descriptionUseCase) GetAllDescriptions(filter map[string]interface{}, page, perPage uint64) (*domain.Descriptions, error) {
	descriptions, err := b.DescriptionRepo.FindAll(filter, page, perPage)
	if err != nil {
		return nil, err
	}
	return descriptions, nil
}

func (b *descriptionUseCase) GetDescription(filter map[string]interface{}) (*domain.Description, error) {
	description, err := b.DescriptionRepo.Find(filter)
	if err != nil {
		return nil, err
	}
	return description, nil
}

func (b *descriptionUseCase) CreateDescription(createDescription domain.Description) (*domain.Description, error) {
	newDescription, err := b.DescriptionRepo.Create(createDescription)
	if err != nil {
		return nil, err
	}
	return newDescription, nil
}

func (b *descriptionUseCase) UpdateDescription(updateBook domain.Description, filter map[string]interface{}) (error) {
	return nil
}

func (b *descriptionUseCase) DeleteDescription(deleteDescription domain.Description) (error) {
	return b.DescriptionRepo.Delete(deleteDescription)
}
