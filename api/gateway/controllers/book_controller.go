package controllers

import (
	"net/http"

	"bookshelf-web-api_gin_clean/api/gateway/repositories"
	"bookshelf-web-api_gin_clean/api/usecases"

	"github.com/gin-gonic/gin"
	"strconv"
	"log"
	"errors"
	"bookshelf-web-api_gin_clean/api/domain"
)

type bookController struct {
	UseCase usecases.BookUseCase
}

type BookController interface {
	GetAllBooks(c *gin.Context)
	GetBook(c *gin.Context)
	CreateBook(c *gin.Context)
	ChangeBookStatus(c *gin.Context)
	DeleteBook(c *gin.Context)
	UpdateBook(c *gin.Context)
}

func NewBookController(dbConnection repositories.DBConnection) BookController {
	repo := repositories.NewBookRepository(dbConnection)
	u := usecases.NewBookUseCase(repo)
	return &bookController{UseCase: u}
}

type BookForm struct {
	Title      string  `json:"title" binding:"required"`
	AuthorID   uint64  `json:"author_id"`
}

type BookUpdateForm struct {
	ID         uint64  `json:"id" binding:"required"`
	Title      string  `json:"title" binding:"required"`
	AuthorID   uint64  `json:"author_id"`
}

type Response struct {
	Content interface{} `json:"content"`
}

func parseStatus(s string) (*domain.ReadState, error) {
	n := domain.NotReadValue
	r := domain.ReadingValue
	r2 := domain.ReadValue
	switch s {
	case "not_read":
		return &n, nil
	case "reading":
		return &r, nil
	case "read":
		return &r2, nil
	default:
		return nil, errors.New("invalid read status")
	}
}

func (b *bookController) GetAllBooks(c *gin.Context) {
	filter := map[string]interface{}{}

	accountId, ok := c.MustGet("account_id").(string)
	if !ok {
		log.Println("GetBook: ", errors.New("accountId parser error"))
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}
	usecases.ByAccountId(filter, accountId)

	page, perPage, err := GetPaginate(c)
	if err != nil {
		log.Println("GetPaginate: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
		return
	}

	sortKey := c.Query("sort_key")

	readStatusStr := c.Query("status")
	if readStatusStr != "" {
		readStatus, err := parseStatus(readStatusStr)
		if err != nil {
			log.Println("GetPaginate: ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusBadRequest})
			return
		}
		usecases.ByStatus(filter, *readStatus)
	}

	books, err := b.UseCase.GetAllBooks(filter, page, perPage, sortKey)
	if err != nil {
		log.Println("GetAllBooks: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, Response{Content: books})
}

func (b *bookController) GetBook(c *gin.Context) {
	bookId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("GetBook: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}
	accountId, ok := c.MustGet("account_id").(string)
	if !ok {
		log.Println("GetBook: ", errors.New("accountId parser error"))
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	filter := usecases.NewFilter()
	usecases.ById(filter, bookId)
	usecases.ByAccountId(filter, accountId)

	book, err := b.UseCase.GetBook(filter)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusOK, Response{Content: book})
}

func (b *bookController) CreateBook(c *gin.Context) {
	form := BookForm{}
	err := c.ShouldBind(&form)
	if err != nil {
		log.Println("CreateBook: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	accountId, ok := c.MustGet("account_id").(string)
	if !ok {
		log.Println("CreateBook: ", errors.New("accountId parser error"))
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}
	book := domain.NewBook()
	book.Title = form.Title
	book.AccountID = accountId
	author := domain.Author{}
	author.ID = form.AuthorID
	book.Author = &author
	
	book.ReadState = domain.NotReadValue

	newBook, err := b.UseCase.CreateBook(book)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusOK, Response{Content: newBook})
}

func (b *bookController) ChangeBookStatus(c *gin.Context) {
	bookId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("ChangeBookStatus: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}
	accountId, ok := c.MustGet("account_id").(string)
	if !ok {
		log.Println("ChangeBookStatus: ", errors.New("accountId parser error"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	filter := usecases.NewFilter()
	usecases.ById(filter, bookId)
	usecases.ByAccountId(filter, accountId)

	err = b.UseCase.ChangeStatus(filter)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.Status(http.StatusOK)
}

func (b *bookController) DeleteBook(c *gin.Context) {
	bookId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("DeleteBook: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}
	accountId, ok := c.MustGet("account_id").(string)
	if !ok {
		log.Println("DeleteBook: ", errors.New("accountId parser error"))
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}
	bookFilter := usecases.NewFilter()
	usecases.ById(bookFilter, bookId)
	usecases.ByAccountId(bookFilter, accountId)

	err = b.UseCase.DeleteBook(bookFilter)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.Status(http.StatusOK)
}

func (b *bookController) UpdateBook(c *gin.Context) {
	form := BookUpdateForm{}
	err := c.ShouldBind(&form)
	if err != nil {
		log.Println("UpdateBook: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}
	accountId, ok := c.MustGet("account_id").(string)
	if !ok {
		log.Println("UpdateBook: ", errors.New("accountId parser error"))
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	filter := usecases.NewFilter()
	usecases.ByAccountId(filter, accountId)
	usecases.ById(filter, form.ID)
	book, err := b.UseCase.GetBook(filter)
	if err != nil {
		log.Println("UpdateBook: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	book.Title = form.Title
	author := domain.Author{}
	author.ID = form.AuthorID
	book.Author = &author

	updatedBook, err := b.UseCase.UpdateBook(*book, nil)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusOK, Response{Content: updatedBook})
}
