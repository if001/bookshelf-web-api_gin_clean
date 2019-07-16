package domain

import (
	"github.com/go-sql-driver/mysql"
	"time"
)

type Book struct {
	Base
	AccountID      string       `json:"account_id"`
	Title          string       `json:"title"`
	Isbn           *string      `json:"isbn"`
	Author         *Author      `json:"author"`
	Publisher      *Publisher   `json:"publisher"`
	StartAt        NullTime     `json:"start_at"`
	EndAt          NullTime     `json:"end_at"`
	ReadState      ReadState    `json:"read_state"`
	Descriptions   Descriptions `json:"descriptions"`
	SmallImageUrl  *string      `json:"small_image_url"`
	MediumImageUrl *string      `json:"medium_image_url"`
	ItemUrl        *string      `json:"item_url"`
	AffiliateUrl   *string      `json:"affiliate_url"`
}

type Books []Book

func NewBook() Book {
	b := Book{}
	b.ID = 0
	b.Title = ""
	b.Isbn = nil
	b.AccountID = ""
	b.Author = nil
	b.Publisher = nil
	b.StartAt = NullTime{mysql.NullTime{Time: time.Now(), Valid: false}}
	b.EndAt = NullTime{mysql.NullTime{Time: time.Now(), Valid: false}}
	b.UpdatedAt = time.Now()
	b.CreatedAt = time.Now()
	b.SmallImageUrl = nil
	b.MediumImageUrl = nil
	return b
}

type PaginateBooks struct {
	Books      Books `json:"books"`
	TotalCount int64 `json:"total_count"`
}

type ReadState int8

const (
	NotReadValue ReadState = iota + 1
	ReadingValue
	ReadValue
)

//func (b *Book) GetReadState() ReadState {
//	if b.StartAt.Valid && b.EndAt.Valid {
//		return ReadValue
//	} else if b.StartAt.Valid && !b.EndAt.Valid {
//		return ReadingValue
//	} else if !b.StartAt.Valid && !b.EndAt.Valid {
//		return NotReadValue
//	} else {
//		return NotReadValue
//	}
//}

func (b *Book) SetStartState() {
	b.StartAt = NullTime{mysql.NullTime{Time: time.Now(), Valid: true}}
	b.EndAt = NullTime{mysql.NullTime{Valid: false}}
	b.ReadState = ReadingValue
}
func (b *Book) SetEndState() {
	b.EndAt = NullTime{mysql.NullTime{Time: time.Now(), Valid: true}}
	b.ReadState = ReadValue
}

type Author struct {
	Base
	Name string `json:"name"`
}

func (Author) TableName() string {
	return "author"
}

type Authors []Author

func (a Authors) FindById(id uint64) *Author {
	for _, v := range a {
		if v.ID == id {
			return &v
		}
	}
	return nil
}

type CountedAuthor struct {
	Author
	Count int64 `json:"count"`
}

type CountedAuthors []CountedAuthor

type Description struct {
	Base
	BookId  uint64 `json:"book_id"`
	Content string `json:"content"`
}

func (Description) TableName() string {
	return "description"
}

type Descriptions []Description

type Publisher struct {
	Base
	Name string `json:"name"`
}

func (Publisher) TableName() string {
	return "publisher"
}

type Publishers []Publisher

func (p Publishers) FindById(id uint64) *Publisher {
	for _, v := range p {
		if v.ID == id {
			return &v
		}
	}
	return nil
}

type CountedPublisher struct {
	Publisher
	Count int64 `json:"count"`
}

type CountedPublishers []CountedPublisher
