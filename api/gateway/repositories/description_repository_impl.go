package repositories

import (
	"bookshelf-web-api_gin_clean/api/usecases"
	"bookshelf-web-api_gin_clean/api/domain"
	"fmt"
	"errors"
)

type DescriptionRepository struct {
	Connection DBConnection
}

func NewDescriptionRepository(conn DBConnection) usecases.DescriptionRepository {
	return &DescriptionRepository{Connection: conn}
}

func (d *DescriptionRepository) FindAll(filter map[string]interface{}, page uint64, perPage uint64) (*domain.Descriptions, error) {
	var descriptions = make(domain.Descriptions, 0)
	query := d.Connection.Where(filter)

	if page > 0 && perPage > 0 {
		query = query.Paginate(page, perPage)
	} else {
		query = d.Connection.Where(filter)
	}

	err := query.SortDesc("updated_at").Bind(&descriptions).HasError()
	if err != nil {
		return nil, fmt.Errorf("FindAll: %s", err)
	}
	return &descriptions, nil
}

func (d *DescriptionRepository) Find(filter map[string]interface{}) (*domain.Description, error) {
	return nil, nil
}

func (d *DescriptionRepository) Create(description domain.Description) (desc *domain.Description, err error) {
	tx := d.Connection.TX()
	defer func() {
		rcv := recover()
		if rcv != nil {
			err := tx.TxRollback()
			if err == nil {
				err = errors.New("in recover: " + rcv.(string))
			}
		}
	}()
	
	err = tx.Create(&description).HasError()
	if err != nil {
		err = fmt.Errorf("description create: %s", err)
		return
	}

	book := domain.Book{}
	filter := map[string]interface{}{"id": description.BookId}
	err = tx.Where(filter).Bind(&book).HasError()
	if err != nil {
		err = fmt.Errorf("description create: %s", err)
		return
	}

	book.UpdatedAt = domain.JstNow()
	err = tx.Update(&book).HasError()
	if err != nil {
		err = fmt.Errorf("description create: %s", err)
		return
	}

	err = tx.TxExec()
	if err != nil {
		err = fmt.Errorf("description create: %s", err)
		return
	}
	return &description, nil
}

func (d *DescriptionRepository) Delete(description domain.Description) error {
	return d.Connection.Delete(description).HasError()
}
