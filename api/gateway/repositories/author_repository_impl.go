package repositories

import (
	"bookshelf-web-api_gin_clean/api/usecases"
	"bookshelf-web-api_gin_clean/api/domain"
	"github.com/jinzhu/gorm"
	"fmt"
)

type AuthorRepository struct {
	Connection DBConnection
}

func NewAuthorRepository(conn DBConnection) usecases.AuthorRepository {
	return &AuthorRepository{Connection: conn}
}

func (a *AuthorRepository) FindAll() (*domain.CountedAuthors, error) {
	author := domain.CountedAuthors{}
	err := a.Connection.CountedAuthorQuery(&author)
	if err != nil {
		return nil, err
	}
	return &author, err
}

func (a *AuthorRepository) Create(author domain.Author) (*domain.Author, error) {
	err := a.Connection.Create(&author).HasError()
	if err != nil {
		return nil, err
	}
	return &author, nil
}

func (a *AuthorRepository) Find(filter map[string]interface{}) (*domain.Author, error) {
	query := a.Connection.Where(filter)

	var author = domain.Author{}
	err := query.Bind(&author).HasError()
	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Find Author: %s", err)
	}
	return &author, nil
}
