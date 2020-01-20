package usecases

import "bookshelf-web-api_gin_clean/api/domain"

type publisherUseCase struct {
	PublisherRepo PublisherRepository
}

type PublisherUseCase interface {
	GetAllPublisher() (*domain.CountedPublishers, error)
	CreatePublisher(createPublisher domain.Publisher) (*domain.Publisher, error)
	GetPublisher(filter map[string]interface{}) (*domain.Publisher, error)
}

func NewPublisherUseCase(publisherRepo PublisherRepository) PublisherUseCase {
	return &publisherUseCase{PublisherRepo: publisherRepo}
}

func (p *publisherUseCase) GetAllPublisher() (*domain.CountedPublishers, error) {
	publishers, err := p.PublisherRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return publishers, nil
}

func (p *publisherUseCase) CreatePublisher(createPublisher domain.Publisher) (*domain.Publisher, error) {
	publisher, err := p.PublisherRepo.Create(createPublisher)
	if err != nil {
		return nil, err
	}
	return publisher, nil
}

func (p *publisherUseCase) GetPublisher(filter map[string]interface{}) (*domain.Publisher, error) {
	return p.PublisherRepo.Find(filter)
}
