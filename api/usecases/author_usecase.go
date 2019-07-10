package usecases

import "bookshelf-web-api_gin_clean/api/domain"

type authorUseCase struct {
	AuthorRepo AuthorRepository
}
type AuthorUseCase interface {
	GetAllAuthor() (*domain.CountedAuthors, error)
}

func NewAuthorUseCase(authorRepo AuthorRepository) AuthorUseCase {
	return &authorUseCase{AuthorRepo: authorRepo}
}

func (a *authorUseCase) GetAllAuthor() (*domain.CountedAuthors, error) {
	authors, err := a.AuthorRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return authors, nil
}