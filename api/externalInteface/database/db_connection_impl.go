package database

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"bookshelf-web-api_gin_clean/api/gateway/repositories"
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

func  (conn *dbConnection) TX() repositories.DBConnection {
	return &dbConnection{DB: conn.DB.Begin()}
}

func (conn *dbConnection) TxRollback() error {
	return conn.DB.Rollback().Error
}

func (conn *dbConnection) TxExec() error {
	return conn.DB.Commit().Error
}

func(conn *dbConnection) CountedAuthorQuery(bind interface{}) error {
	count := 0
	return conn.DB.Table("books").
		Joins("left join author on books.author_id = author.id").
		Group("author.id").
		Select("author.id, author.name, author.created_at, author.updated_at, count(*) as count").
		Count(&count).
		Having("count >? ",0).
		Having("author.id is not NULL").
		Find(bind).
		Error
}

func (conn *dbConnection) HasError() error {
	return conn.DB.Error
}

func NewSqlConnection() dbConnection {
	config, err := LoadConfig()
	if err != nil {
		panic(err.Error())
	}

	dbconf := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DB.User,
		config.DB.Password,
		config.DB.Host,
		config.DB.DB)
	db, err := gorm.Open("mysql", dbconf)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)

	// defer db.Close()
	return dbConnection{DB: db}
}
