package repositories

import (
	"bookshelf-web-api_gin_clean/api/usecases"
	"bookshelf-web-api_gin_clean/api/domain"
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
