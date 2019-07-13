package repositories

import (
	"bookshelf-web-api_gin_clean/api/domain"
	"bookshelf-web-api_gin_clean/api/usecases"
)

type PublisherRepository struct {
	Connection DBConnection
}

func NewPublisherRepository(conn DBConnection) usecases.PublisherRepository {
	return &PublisherRepository{Connection: conn}
}

func (p *PublisherRepository) FindAll() (*domain.CountedPublishers, error) {
	publishers := domain.CountedPublishers{}
	err := p.Connection.CountedPublisherQuery(&publishers)
	if err != nil {
		return nil, err
	}
	return &publishers, err
}

func (p *PublisherRepository) Create(publisher domain.Publisher) (*domain.Publisher, error) {
	err := p.Connection.Create(&publisher).HasError()
	if err != nil {
		return nil, err
	}
	return &publisher, nil
}
