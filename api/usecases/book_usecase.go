package usecases

import (
	"bookshelf-web-api_gin_clean/api/domain"
	"errors"
)

type bookUseCase struct {
	BookRepo BookRepository
}
type BookUseCase interface {
	GetAllBooks(filter map[string]interface{}, page, parPage uint64, sortKey string) (*domain.PaginateBooks, error) // TODO paging
	GetBook(filter map[string]interface{}) (*domain.Book, error)
	UpdateBook(updateBook domain.Book, filter map[string]interface{}) (*domain.Book, error)
	CreateBook(createBook domain.Book) (*domain.Book, error)
	DeleteBook(filter map[string]interface{}) error

	ChangeStatus(filter map[string]interface{}) error
	// StoreCategories() error
	// ChangeRating() error
	// SetNextBook() error
	// SetPrevBook() error
}

func NewBookUseCase(repo BookRepository) BookUseCase {
	return &bookUseCase{BookRepo: repo}
}

func (b *bookUseCase) GetAllBooks(filter map[string]interface{}, page, perPage uint64, sortKey string) (*domain.PaginateBooks, error) {
	books, err := b.BookRepo.FindAll(filter, page, perPage, sortKey)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (b *bookUseCase) GetBook(filter map[string]interface{}) (*domain.Book, error) {
	book, err := b.BookRepo.Find(filter)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (b *bookUseCase) UpdateBook(updateBook domain.Book, filter map[string]interface{}) (*domain.Book, error) {
	err := b.BookRepo.Store(updateBook)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (b *bookUseCase) CreateBook(createBook domain.Book) (*domain.Book, error) {
	newBook, err := b.BookRepo.Create(createBook)
	if err != nil {
		return nil, err
	}
	return newBook, nil
}

func (b *bookUseCase) DeleteBook(filter map[string]interface{}) (error) {
	// TODO 関連するdescription、カテゴリなどの削除
	err := b.BookRepo.Delete(filter)
	if err != nil {
		return err
	}
	return nil
}

func (b *bookUseCase) ChangeStatus(filter map[string]interface{}) (error) {
	book, err := b.BookRepo.Find(filter)
	if err != nil {
		return err
	}

	switch book.ReadState {
	case domain.NotReadValue:
		book.SetStartState()
	case domain.ReadingValue:
		book.SetEndState()
	case domain.ReadValue:
		book.SetStartState()
	default:
		return errors.New("ChangeStatus: bad status")
	}
	err = b.BookRepo.Store(*book)
	if err != nil {
		return err
	}
	return nil
}
