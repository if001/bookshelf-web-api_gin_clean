package repositories

import (
	"bookshelf-web-api_gin_clean/api/domain"
	"bookshelf-web-api_gin_clean/api/usecases"
	"time"
	"fmt"
	"errors"
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
	AuthorID       *uint64
	PublisherID    *uint64
	StartAt        domain.NullTime
	EndAt          domain.NullTime
	ReadState      domain.ReadState
	SmallImageUrl  *string
	MediumImageUrl *string
}

func (BookTable) TableName() string {
	return "books"
}
func (b *BookTable) ToModel() domain.Book {
	m := domain.Book{
		AccountID: b.AccountID,
		Title:     b.Title,
		Author:    nil,
		Publisher: nil,
		StartAt:   b.StartAt,
		EndAt:     b.EndAt,
		ReadState: b.ReadState,
	}
	m.ID = b.ID

	m.SmallImageUrl = b.SmallImageUrl
	m.MediumImageUrl = b.MediumImageUrl
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
		Title:       b.Title,
		AccountID:   b.AccountID,
		AuthorID:    authorID,
		PublisherID: publisherID,
		StartAt:     b.StartAt,
		EndAt:       b.EndAt,
		ReadState:   b.ReadState,
	}
	t.ID = b.ID
	t.SmallImageUrl = b.SmallImageUrl
	t.MediumImageUrl = b.MediumImageUrl
	t.UpdatedAt = b.UpdatedAt
	t.CreatedAt = b.CreatedAt
	return t
}

func NewBookRepository(conn DBConnection) usecases.BookRepository {
	return &BookRepository{Connection: conn}
}

func (b *BookRepository) FindAll(filter map[string]interface{}, page uint64, perPage uint64, sortKey string) (*domain.PaginateBooks, error) {
	var bookTables = make([]BookTable, 0)
	var count int64 = 0
	if err := b.Connection.Table(&bookTables).Select(filter).Count(&count).HasError(); err != nil {
		return nil, err
	}

	query := b.Connection.Select(filter)
	if page > 0 && perPage > 0 {
		query = query.Paginate(page, perPage)
	}
	if sortKey == "" {
		query = query.SortDesc("updated_at")
	} else if sortKey == "title" {
		query = query.SortAsc(sortKey)
	} else {
		query = query.SortDesc(sortKey)
	}

	err := query.Bind(&bookTables).HasError()
	if err != nil {
		return nil, fmt.Errorf("FindAll: %s", err)
	}
	authorTables := domain.Authors{}
	cc := b.Connection
	for _, v := range bookTables {
		cc.OrFilter(map[string]interface{}{"author_id": v.AuthorID})
	}
	err = cc.Bind(&authorTables).HasError()
	if err != nil {
		return nil, fmt.Errorf("FindAll: %s", err)
	}

	publisherTables := domain.Publishers{}
	pc := b.Connection
	for _, v := range bookTables {
		pc.OrFilter(map[string]interface{}{"publisher_id": v.PublisherID})
	}
	err = cc.Bind(&publisherTables).HasError()
	if err != nil {
		return nil, fmt.Errorf("FindAll: %s", err)
	}


	books := domain.Books{}
	for _, v := range bookTables {
		b := v.ToModel()
		if v.AuthorID != nil {
			author := authorTables.FindById(*v.AuthorID)
			b.Author = author
		} else {
			b.Author = nil
		}
		if v.PublisherID != nil {
			publisher := publisherTables.FindById(*v.PublisherID)
			b.Publisher = publisher
		} else {
			b.Publisher = nil
		}
		books = append(books, b)
	}

	paginateBooks := domain.PaginateBooks{
		Books:      books,
		TotalCount: count,
	}

	return &paginateBooks, nil
}

func (b *BookRepository) Find(filter map[string]interface{}) (*domain.Book, error) {
	var bookTable = BookTable{}
	err := b.Connection.Select(filter).Bind(&bookTable).HasError()
	if err != nil {
		return nil, err
	}

	book := bookTable.ToModel()
	if bookTable.AuthorID == nil {
		return &book, nil
	}
	authorTable := domain.Author{}
	authorFilter := map[string]interface{}{"id": bookTable.AuthorID}
	err = b.Connection.Select(authorFilter).Bind(&authorTable).HasError()
	if err != nil {
		return nil, err
	}
	book.Author = &authorTable

	if bookTable.PublisherID == nil {
		return &book, nil
	}
	publisherTable := domain.Publisher{}
	publisherFilter := map[string]interface{}{"id": bookTable.PublisherID}
	err = b.Connection.Select(publisherFilter).Bind(&publisherTable).HasError()
	if err != nil {
		return nil, err
	}
	book.Publisher = &publisherTable

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
	err = b.Connection.Select(filter).Bind(&bookTable).HasError()
	if err != nil {
		return
	}
	var descriptionsTable = domain.Descriptions{}
	m := map[string]interface{}{"book_id": bookTable.ID}
	err = b.Connection.Select(m).Bind(&descriptionsTable).HasError()
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
	t.UpdatedAt = time.Now()
	return b.Connection.Update(t).HasError()
}

func (b *BookRepository) UpdateUpdatedAt(filter map[string]interface{}) error {
	var bookTable = BookTable{}
	err := b.Connection.Select(filter).Bind(&bookTable).HasError()
	if err != nil {
		return err
	}
	if bookTable.ID == 0 {
		return errors.New("Store:TableNotFound")
	}

	bookTable.UpdatedAt = time.Now()
	return b.Connection.Update(bookTable).HasError()
}
