package database

import (
	"bookshelf-web-api_gin_clean/api/gateway/repositories"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type dbConnection struct {
	DB *gorm.DB
}

func (conn *dbConnection) Bind(bind interface{}) repositories.DBConnection {
	return &dbConnection{DB: conn.DB.Find(bind)}
}

func (conn *dbConnection) Paginate(page, perPage uint64) repositories.DBConnection {
	return &dbConnection{DB: conn.DB.Offset(perPage * (page - 1)).Limit(perPage)}
}

func (conn *dbConnection) Select(filter interface{}) repositories.DBConnection {
	return &dbConnection{DB: conn.DB.Where(filter)}
}

func (conn *dbConnection) OrFilter(filter interface{}) repositories.DBConnection {
	return &dbConnection{DB: conn.DB.Or(filter)}
}

func (conn *dbConnection) Create(data interface{}) repositories.DBConnection {
	return &dbConnection{DB: conn.DB.Create(data)}
}

func (conn *dbConnection) Update(data interface{}) repositories.DBConnection {
	return &dbConnection{DB: conn.DB.Save(data)}
}

func (conn *dbConnection) Delete(data interface{}) repositories.DBConnection {
	return &dbConnection{DB: conn.DB.Delete(data)}
}

func (conn *dbConnection) SortDesc(key string) repositories.DBConnection {
	return &dbConnection{DB: conn.DB.Order(fmt.Sprintf("%s desc", key), true)}
}

func (conn *dbConnection) SortAsc(key string) repositories.DBConnection {
	return &dbConnection{DB: conn.DB.Order(fmt.Sprintf("%s asc", key), true)}
}

func (conn *dbConnection) Count(count *int64) repositories.DBConnection {
	return &dbConnection{DB: conn.DB.Count(count)}
}

func (conn *dbConnection) Table(table interface{}) repositories.DBConnection {
	return &dbConnection{DB: conn.DB.Model(table)}
}

func (conn *dbConnection) TX() repositories.DBConnection {
	return &dbConnection{DB: conn.DB.Begin()}
}

func (conn *dbConnection) TxRollback() error {
	return conn.DB.Rollback().Error
}

func (conn *dbConnection) TxExec() error {
	return conn.DB.Commit().Error
}

func (conn *dbConnection) CountedAuthorQuery(bind interface{}) error {
	count := 0
	return conn.DB.Table("author").
		Joins("left join books on books.author_id = author.id").
		Group("author.id").
		Select("author.id, author.name, author.created_at, author.updated_at, count(*) as count").
		Count(&count).
		Having("count >? ", 0).
		Having("author.id is not NULL").
		Find(bind).
		Error
}

func (conn *dbConnection) CountedPublisherQuery(bind interface{}) error {
	count := 0
	return conn.DB.Table("publisher").
		Joins("left join books on books.publisher_id = publisher.id").
		Group("publisher.id").
		Select("publisher.id, publisher.name, publisher.created_at, publisher.updated_at, count(*) as count").
		Count(&count).
		Having("count >? ", 0).
		Having("publisher.id is not NULL").
		Find(bind).
		Error
}

func (conn *dbConnection) SelectBookWith(bind interface{}) repositories.DBConnection {
	return &dbConnection{DB: conn.DB.Table("books").
		Select("books.*, author.id, author.name,author.created_at,author.updated_at, " +
			"publisher.id, publisher.name, publisher.created_at, publisher.updated_at").
		Joins("left join author on author.id = books.author_id").
		Joins("left join publisher on publisher.id = books.publisher_id").
		Find(bind)}
}

func (conn *dbConnection) HasError() error {
	return conn.DB.Error
}

func NewSqlConnection(url string) dbConnection {
	db, err := gorm.Open("mysql", url)
	if err != nil {
		panic(err.Error())
		// log.Errorf(ctx, "gormOpen: %s", err)
	}
	db.LogMode(true)

	// defer db.Close()
	return dbConnection{DB: db}
}
