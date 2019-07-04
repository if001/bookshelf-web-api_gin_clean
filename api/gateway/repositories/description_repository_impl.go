package repositories


import (
	"bookshelf-web-api_gin_clean/api/usecases"
	"bookshelf-web-api_gin_clean/api/domain"
	"fmt"
)

type DescriptionRepository struct {
	Connection DBConnection
}

func NewDescriptionRepository(conn DBConnection) usecases.DescriptionRepository {
	return &DescriptionRepository{Connection: conn}
}

func (d *DescriptionRepository) FindAll(filter map[string]interface{}, page uint64, perPage uint64) (*domain.Descriptions, error) {
	var descriptions = make(domain.Descriptions, 0)
	if page > 0 && perPage > 0 {
		err := d.Connection.Select(filter).Paginate(page, perPage).Bind(&descriptions).HasError()
		if err != nil {
			return nil, fmt.Errorf("FindAll: %s",err)
		}
	} else {
		err := d.Connection.Select(filter).Bind(&descriptions).HasError()
		if err != nil {
			return nil, fmt.Errorf("FindAll: %s",err)
		}
	}
	return &descriptions, nil
}

func (d *DescriptionRepository) Find(filter map[string]interface{}) (*domain.Description, error) {
	return nil, nil
}

func (d *DescriptionRepository) Create(description domain.Description) (*domain.Description, error) {
	err := d.Connection.Create(&description).HasError()
	if err != nil {
		return nil, fmt.Errorf("description create: %s",err)
	}
	return &description, nil
}

func (d *DescriptionRepository) Delete(description domain.Description) error {
	return d.Connection.Delete(description).HasError()
}
