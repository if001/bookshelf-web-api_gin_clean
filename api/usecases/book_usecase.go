package usecases

import (
	"bookshelf-web-api_gin_clean/api/domain"
	"errors"
	"strings"
	"unicode"
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
	CountByName(filter map[string]interface{}, key string) (*domain.CountedNames, error)
	CountByDate(filter map[string]interface{}, dateKey , dateType string) (*domain.CountedDates, error)
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

func containsName(s []domain.CountedName, name string) bool {
	for _, v := range s {
		if v.Name == name {
			return true
		}
	}
	return false
}

func spaceMap(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}

func (b *bookUseCase) CountByName(filter map[string]interface{}, key string) (*domain.CountedNames, error) {
	var err error = nil
	tmp := make(domain.CountedNames, 0)
	bookWithName := &tmp
	if key == "publisher" {
		bookWithName, err = b.BookRepo.CountByPublisher(filter)
	} else if key == "author" {
		bookWithName, err = b.BookRepo.CountByAuthor(filter)
	} else {
		return nil, errors.New("CountByName: invalid key")
	}
	if err != nil {
		return nil, err
	}

	tmp2 := bookWithName.RemoveByName("")
	bookWithName = &tmp2

	length := len(*bookWithName)
	bookWithNameRe := make(domain.CountedNames, length)
	for i, v := range *bookWithName {
		if strings.Contains(v.Name, "/") {
			authorExTranslation := strings.Split(v.Name, "/")[0]
			bookWithNameRe[i].Name = spaceMap(authorExTranslation)
			bookWithNameRe[i].Count = v.Count
		} else {
			bookWithNameRe[i].Name = spaceMap(v.Name)
			bookWithNameRe[i].Count = v.Count
		}
	}

	bookWithNameGroupBy := make(domain.CountedNames, 0)
	for _, v := range bookWithNameRe {
		if !containsName(bookWithNameGroupBy, v.Name) {
			bookWithNameGroupBy = append(
				bookWithNameGroupBy,
				domain.CountedName{
					Name:  v.Name,
					Count: v.Count,
				})
		} else {
			if index, ok := bookWithNameGroupBy.SearchIndex(v.Name); ok {
				bookWithNameGroupBy[index].Count += v.Count
			}
		}
	}
	authorCountedNameGroupBySliced := bookWithNameGroupBy.SortByCount()
	if len(authorCountedNameGroupBySliced) >= 20 {
		s := authorCountedNameGroupBySliced[:20]
		return &s, nil
	} else {
		return &authorCountedNameGroupBySliced, nil
	}
}

func (b *bookUseCase) CountByDate(filter map[string]interface{}, dateKey , dateType string) (*domain.CountedDates, error) {
	format := ""
	if dateType == domain.DateKeyDaily {
		format = "%Y-%m-%d"
	} else if dateType == domain.DateKeyMonthly {
		format = "%Y-%m"
	} else {
		return nil, errors.New("CountByDate: invalid date key")
	}

	countedDates, err := b.BookRepo.CountByDate(filter, dateKey, format)
	if err != nil {
		return nil, err
	}
	return countedDates, nil
}
