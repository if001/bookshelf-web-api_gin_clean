package repositories

import (
	"bookshelf-web-api_gin_clean/api/domain"
	"bookshelf-web-api_gin_clean/api/usecases"
	"errors"
	"fmt"
	"time"
)

type BookRepository struct {
	Connection DBConnection
}

type Base struct {
	ID        uint64    `gorm:"primary_key" sql:"AUTO_INCREMENT"`
	CreatedAt time.Time `sql:"not null;type:date"`
	UpdatedAt time.Time `sql:"not null;type:date"`
}
type BookTable struct {
	Base
	Title          string
	AccountID      string
	Isbn           *string
	AuthorID       *uint64
	PublisherID    *uint64
	StartAt        domain.NullTime
	EndAt          domain.NullTime
	ReadState      domain.ReadState
	SmallImageUrl  *string
	MediumImageUrl *string
	ItemUrl        *string
	AffiliateUrl   *string
}

func (BookTable) TableName() string {
	return "books"
}
func (b *BookTable) ToModel() domain.Book {
	m := domain.Book{
		AccountID:      b.AccountID,
		Title:          b.Title,
		Isbn:           b.Isbn,
		Author:         nil,
		Publisher:      nil,
		StartAt:        b.StartAt,
		EndAt:          b.EndAt,
		ReadState:      b.ReadState,
		SmallImageUrl:  b.SmallImageUrl,
		MediumImageUrl: b.MediumImageUrl,
		ItemUrl:        b.ItemUrl,
		AffiliateUrl:   b.AffiliateUrl,
	}
	m.ID = b.ID
	m.CreatedAt = b.CreatedAt
	m.UpdatedAt = b.UpdatedAt
	return m
}
func ToTable(b domain.Book) BookTable {
	var authorID *uint64 = nil
	if b.Author != nil {
		authorID = &b.Author.ID
	}
	var publisherID *uint64 = nil
	if b.Publisher != nil {
		publisherID = &b.Publisher.ID
	}

	t := BookTable{
		Title:          b.Title,
		AccountID:      b.AccountID,
		Isbn:           b.Isbn,
		AuthorID:       authorID,
		PublisherID:    publisherID,
		StartAt:        b.StartAt,
		EndAt:          b.EndAt,
		ReadState:      b.ReadState,
		SmallImageUrl:  b.SmallImageUrl,
		MediumImageUrl: b.MediumImageUrl,
		ItemUrl:        b.ItemUrl,
		AffiliateUrl:   b.AffiliateUrl,
	}
	t.ID = b.ID
	t.UpdatedAt = b.UpdatedAt
	t.CreatedAt = b.CreatedAt
	return t
}

type BookWith struct {
	domain.Book
	*domain.Author
	*domain.Publisher
}

func NewBookRepository(conn DBConnection) usecases.BookRepository {
	return &BookRepository{Connection: conn}
}

func (b *BookRepository) FindAll(filter map[string]interface{}, page uint64, perPage uint64, sortKey string) (*domain.PaginateBooks, error) {
	query := b.Connection

	if bookFilter, ok := filter["book"]; ok {
		bookFilterStr, ok := bookFilter.(string)
		if !ok {
			return nil, errors.New("FindAll Book Filter Error")
		}
		query = query.SearchBook(bookFilterStr)

		delete(filter, "book")
	}
	query = query.Where(filter)

	queryForCount := DBConnection(query)
	var bookWiths = make([]BookWith, 0)
	var count int64 = 0
	err := queryForCount.SelectBookWith().Bind(&bookWiths).Count(&count).HasError()
	if err != nil {
		return nil, fmt.Errorf("FindAll: %s", err)
	}

	if page > 0 && perPage > 0 {
		query = query.Paginate(page, perPage)
	}
	if sortKey == "" {
		query = query.SortDesc("books.updated_at")
	} else if sortKey == "title" {
		query = query.SortAsc("books." + sortKey)
	} else {
		query = query.SortDesc("books." + sortKey)
	}

	err = query.SelectBookWith().Bind(&bookWiths).HasError()
	if err != nil {
		return nil, fmt.Errorf("FindAll: %s", err)
	}
	var books = make(domain.Books, 0)
	for _, v := range bookWiths {
		book := domain.Book{}
		book = v.Book
		if v.Author.ID != 0 {
			book.Author = v.Author
		} else {
			book.Author = nil
		}
		if v.Publisher.ID != 0 {
			book.Publisher = v.Publisher
		} else {
			book.Publisher = nil
		}
		books = append(books, book)
	}

	paginateBooks := domain.PaginateBooks{
		Books:      books,
		TotalCount: count,
	}

	return &paginateBooks, nil
}

func (b *BookRepository) Find(filter map[string]interface{}) (*domain.Book, error) {
	query := b.Connection.Where(filter)
	var bookWith = BookWith{}
	err := query.SelectBookWith().Bind(&bookWith).HasError()
	if err != nil {
		return nil, fmt.Errorf("FindAll: %s", err)
	}

	var book = domain.Book{}
	book = bookWith.Book
	// 初期化された構造体が入るので、nilを入れ直している
	// TODO 修正したい
	if bookWith.Author.ID != 0 {
		book.Author = bookWith.Author
	} else {
		book.Author = nil
	}
	if bookWith.Publisher.ID != 0 {
		book.Publisher = bookWith.Publisher
	} else {
		book.Publisher = nil
	}
	return &book, nil
}

func (b *BookRepository) Create(book domain.Book) (*domain.Book, error) {
	t := ToTable(book)
	err := b.Connection.Create(&t).HasError()
	if err != nil {
		return nil, err
	}
	newBook := t.ToModel()
	return &newBook, nil
}

func (b *BookRepository) Delete(filter map[string]interface{}) (err error) {
	var bookTable = BookTable{}
	err = b.Connection.Where(filter).Bind(&bookTable).HasError()
	if err != nil {
		return
	}
	var descriptionsTable = domain.Descriptions{}
	m := map[string]interface{}{"book_id": bookTable.ID}
	err = b.Connection.Where(m).Bind(&descriptionsTable).HasError()
	if err != nil {
		return err
	}
	tx := b.Connection.TX()
	defer func() {
		rcv := recover()
		if rcv != nil {
			err = tx.TxRollback()
			if err == nil {
				err = errors.New("in recover: " + rcv.(string))
			}
		}
	}()

	for _, v := range descriptionsTable {
		err = tx.Delete(v).HasError()
		if err != nil {
			err = tx.TxRollback()
			return
		}
	}

	err = tx.Delete(bookTable).HasError()
	if err != nil {
		err = tx.TxRollback()
		return
	}

	err = tx.TxExec()
	if err != nil {
		err = tx.TxRollback()
		return
	}

	return nil
}

func (b *BookRepository) Store(book domain.Book) error {
	t := ToTable(book)
	t.UpdatedAt = domain.JstNow()
	return b.Connection.Update(t).HasError()
}

func (b *BookRepository) UpdateUpdatedAt(filter map[string]interface{}) error {
	var bookTable = BookTable{}
	err := b.Connection.Where(filter).Bind(&bookTable).HasError()
	if err != nil {
		return err
	}
	if bookTable.ID == 0 {
		return errors.New("Store:TableNotFound")
	}

	bookTable.UpdatedAt = domain.JstNow()
	return b.Connection.Update(bookTable).HasError()
}

const (
	countedAuthorKey    = "author.name"
	countedPublisherKey = "publisher.name"
)

func (b *BookRepository) CountByAuthor(filter map[string]interface{}) (*domain.CountedNames, error) {
	return countedBy(b, filter, countedAuthorKey)
}

func (b *BookRepository) CountByPublisher(filter map[string]interface{}) (*domain.CountedNames, error) {
	return countedBy(b, filter, countedPublisherKey)
}

func countedBy(b *BookRepository, filter map[string]interface{}, key string) (*domain.CountedNames, error) {
	query := b.Connection
	if key == countedAuthorKey {
		query = query.Where(filter).SelectBookWithAuthorName()
	} else if key == countedPublisherKey {
		query = query.SelectBookWithPublisherName()
	} else {
		return nil, errors.New("countedBy: invalid counted name key")
	}

	var countedNames = make(domain.CountedNames, 0)
	err := query.GroupBy(key).Bind(&countedNames).HasError()
	if err != nil {
		return nil, fmt.Errorf("countedBy: %s", err)
	}
	return &countedNames, nil
}

func (b *BookRepository) CountByDate(filter map[string]interface{}, key, format string) (*domain.CountedDates, error) {
	var countedDates = make(domain.CountedDates, 0)
	query := b.Connection.Where(filter)

	if key == domain.BookRegister {
		query = query.GroupByDate("created_at", format)
	} else if key == domain.BookReadStart {
		query = query.
			RowWhere(fmt.Sprintf("%s is not null", key)).
			GroupByDate(key, format)
	} else if key == domain.BookReadEnd {
		query = query.
			RowWhere(fmt.Sprintf("%s is not null", key)).
			RowWhere("start_at is not null").
			GroupByDate(key, format)
	} else {
		return nil, errors.New("CountByDate: invalid key")
	}

	// TODO 月指定で取ってこれるようにするまで、limitつけとく
	err := query.Limit(500).Bind(&countedDates).HasError()
	if err != nil {
		return nil, fmt.Errorf("CountByDate: %s", err)
	}
	return &countedDates, nil
}
