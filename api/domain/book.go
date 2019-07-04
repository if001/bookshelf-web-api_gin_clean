package domain

import (
	"github.com/go-sql-driver/mysql"
	"time"
)

type Book struct {
	Base
	AccountID    string       `json:"account_id"`
	Title        string       `json:"title"`
	Author       *Author      `json:"author"`
	StartAt      NullTime     `json:"start_at"`
	EndAt        NullTime     `json:"end_at"`
	ReadState    ReadState    `json:"read_state"`
	Descriptions Descriptions `json:"descriptions"`
}

type Books []Book

func NewBook() Book {
	b := Book{}
	b.ID = 0
	b.Title = ""
	b.AccountID = ""
	b.Author = nil
	b.StartAt = NullTime{mysql.NullTime{Time: time.Now(), Valid: false}}
	b.EndAt = NullTime{mysql.NullTime{Time: time.Now(), Valid: false}}
	b.UpdatedAt = time.Now()
	b.CreatedAt = time.Now()
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

type Description struct {
	Base
	BookId  uint64 `json:"book_id"`
	Content string `json:"content"`
}

func (Description) TableName() string {
	return "description"
}

type Descriptions []Description
