package repositories

import (
	"bookshelf-web-api_gin_clean/api/domain"
	"bookshelf-web-api_gin_clean/api/usecases"
	"github.com/jinzhu/gorm"
	"fmt"
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

func (p *PublisherRepository) Find(filter map[string]interface{}) (*domain.Publisher, error) {
	query := p.Connection.Where(filter)

	var publisher = domain.Publisher{}
	err := query.Bind(&publisher).HasError()
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Find Publisher: %s", err)
	}
	return &publisher, nil
}
