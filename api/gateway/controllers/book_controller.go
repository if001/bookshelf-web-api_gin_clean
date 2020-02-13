package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"bookshelf-web-api_gin_clean/api/domain"
	"bookshelf-web-api_gin_clean/api/gateway/repositories"
	"bookshelf-web-api_gin_clean/api/usecases"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

type bookController struct {
	UseCase          usecases.BookUseCase
	AuthorUseCase    usecases.AuthorUseCase
	PublisherUseCase usecases.PublisherUseCase
}

type BookController interface {
	GetAllBooks(c *gin.Context)
	GetBook(c *gin.Context)
	GetBookLimit(c *gin.Context)
	CreateBook(c *gin.Context)
	CreateBookWith(c *gin.Context)
	ChangeBookStatus(c *gin.Context)
	DeleteBook(c *gin.Context)
	UpdateBook(c *gin.Context)
	GetCountedByAuthor(c *gin.Context)
	GetCountedByPublisher(c *gin.Context)
	GetCountedRegisterDaily(c *gin.Context)
	GetCountedStartDaily(c *gin.Context)
	GetCountedEndDaily(c *gin.Context)
	GetCountedMonthly(c *gin.Context)
}

func NewBookController(dbConnection repositories.DBConnection) BookController {
	repo := repositories.NewBookRepository(dbConnection)
	u := usecases.NewBookUseCase(repo)
	return &bookController{UseCase: u}
}

type BookForm struct {
	Title          string  `json:"title" binding:"required"`
	Isbn           *string `json:"isbn"`
	AuthorID       *uint64 `json:"author_id"`
	PublisherID    *uint64 `json:"publisher_id"`
	SmallImageUrl  *string `json:"small_image_url"`
	MediumImageUrl *string `json:"medium_image_url"`
	ItemUrl        *string `json:"item_url"`
	AffiliateUrl   *string `json:"affiliate_url"`
}

type BookFormWith struct {
	Title          string  `json:"title" binding:"required"`
	Isbn           *string `json:"isbn"`
	AuthorName     *string `json:"author_name"`
	PublisherName  *string `json:"publisher_name"`
	SmallImageUrl  *string `json:"small_image_url"`
	MediumImageUrl *string `json:"medium_image_url"`
	ItemUrl        *string `json:"item_url"`
	AffiliateUrl   *string `json:"affiliate_url"`
}

type BookUpdateForm struct {
	ID          uint64     `json:"id" binding:"required"`
	Title       string     `json:"title" binding:"required"`
	AuthorID    *uint64    `json:"author_id"`
	PublisherID *uint64    `json:"publisher_id"`
	StartAt     *time.Time `json:"start_at"`
	EndAt       *time.Time `json:"end_id"`
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

	isbn := c.Query("isbn")
	if isbn != "" {
		usecases.ByISBN(filter, isbn)
	}

	filter, err := addAccountToFilter(c, &filter)
	if err != nil {
		badRequestWithSentry(c, "GetAllBooks: ", err)
		return
	}

	page, perPage, err := GetPaginate(c)
	if err != nil {
		badRequestWithSentry(c, "GetAllBooks: ", err)
		return
	}

	sortKey := c.Query("sort_key")

	readStatusStr := c.Query("status")
	if readStatusStr != "" {
		readStatus, err := parseStatus(readStatusStr)
		if err != nil {
			badRequestWithSentry(c, "GetAllBooks: ", err)
			return
		}
		usecases.ByStatus(filter, *readStatus)
	}

	bookFilter := c.Query("book")
	if bookFilter != "" {
		usecases.ByBook(filter, bookFilter)
	}

	books, err := b.UseCase.GetAllBooks(filter, page, perPage, sortKey)
	if err != nil {
		internalServerErrorWithSentry(c, "GetAllBooks: ", err)
		return
	}
	c.JSON(http.StatusOK, Response{Content: books})
}

func (b *bookController) GetBook(c *gin.Context) {
	bookId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		badRequestWithSentry(c, "GetBook: ", err)
		return
	}

	filter := usecases.NewFilter()

	filter, err = addAccountToFilter(c, &filter)
	if err != nil {
		badRequestWithSentry(c, "GetBook: ", err)
		return
	}

	usecases.ById(filter, bookId)

	book, err := b.UseCase.GetBook(filter)
	if err != nil {
		internalServerErrorWithSentry(c, "GetBook: ", err)
		return
	}
	c.JSON(http.StatusOK, Response{Content: book})
}

func (b *bookController) GetBookLimit(c *gin.Context) {
	type LimitBook struct {
		Title          string              `json:"title"`
		Isbn           *string             `json:"isbn"`
		Author         *domain.Author      `json:"author"`
		Publisher      *domain.Publisher   `json:"publisher"`
		Descriptions   domain.Descriptions `json:"descriptions"`
		SmallImageUrl  *string             `json:"small_image_url"`
		MediumImageUrl *string             `json:"medium_image_url"`
		ItemUrl        *string             `json:"item_url"`
		AffiliateUrl   *string             `json:"affiliate_url"`
	}

	bookId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		badRequestWithSentry(c, "GetBook: ", err)
		return
	}

	filter := usecases.NewFilter()

	usecases.ById(filter, bookId)

	book, err := b.UseCase.GetBook(filter)
	if err != nil {
		internalServerErrorWithSentry(c, "GetBook: ", err)
		return
	}

	limitBook := LimitBook{
		book.Title,
		book.Isbn,
		book.Author,
		book.Publisher,
		book.Descriptions,
		book.SmallImageUrl,
		book.MediumImageUrl,
		book.ItemUrl,
		book.AffiliateUrl,
	}

	c.JSON(http.StatusOK, Response{Content: limitBook})
}

func (b *bookController) CreateBook(c *gin.Context) {
	form := BookForm{}
	err := c.ShouldBind(&form)
	if err != nil {
		badRequestWithSentry(c, "CreateBook: ", err)
		return
	}

	accountId, ok := c.MustGet("account_id").(string)
	if !ok {
		internalServerErrorWithSentry(c, "CreateBook: ", errors.New("accountId parser error"))
		return
	}

	book := domain.NewBook()
	book.Title = form.Title
	book.AccountID = accountId

	author := domain.Author{}
	if form.AuthorID != nil {
		author.ID = *form.AuthorID
		book.Author = &author
	} else {
		book.Author = nil
	}

	publisher := domain.Publisher{}
	if form.PublisherID != nil {
		publisher.ID = *form.PublisherID
		book.Publisher = &publisher
	} else {
		book.Publisher = nil
	}
	book.Isbn = form.Isbn
	book.SmallImageUrl = form.SmallImageUrl
	book.MediumImageUrl = form.MediumImageUrl
	book.ItemUrl = form.ItemUrl
	book.AffiliateUrl = form.AffiliateUrl
	book.ReadState = domain.NotReadValue

	newBook, err := b.UseCase.CreateBook(book)
	if err != nil {
		log.Println(err.Error())
		if gin.Mode() == gin.ReleaseMode {
			sentryLogError("CreateBook: ", err)
		}
		c.JSON(http.StatusBadRequest, Response{Content: book})
		return
	}
	c.JSON(http.StatusOK, Response{Content: newBook})
}

func (b *bookController) CreateBookWith(c *gin.Context) {
	form := BookFormWith{}
	err := c.ShouldBind(&form)
	if err != nil {
		badRequestWithSentry(c, "CreateBook: ", err)
		return
	}

	accountId, ok := c.MustGet("account_id").(string)
	if !ok {
		internalServerErrorWithSentry(c, "CreateBook: ", errors.New("accountId parser error"))
		return
	}

	book := domain.NewBook()
	book.Title = form.Title
	book.AccountID = accountId

	if form.AuthorName != nil {
		aFilter := usecases.NewFilter()
		usecases.ByName(aFilter, *form.AuthorName)
		author, err := b.AuthorUseCase.GetAuthor(aFilter)
		if err != nil {
			internalServerErrorWithSentry(c, "GetAuthor: ", err)
		}
		book.Author = author
		if author == nil {
			createAuthor := domain.Author{}
			createAuthor.Name = *form.AuthorName
			newAuthor, err := b.AuthorUseCase.CreateAuthor(createAuthor)
			if err != nil {
				log.Println("CreateAuthor: ", err.Error())
				internalServerErrorWithSentry(c, "CreateAuthor: ", err)
				return
			}
			book.Author = newAuthor
		}
	}

	if form.PublisherName != nil {
		pFilter := usecases.NewFilter()
		usecases.ByName(pFilter, *form.PublisherName)
		publisher, err := b.PublisherUseCase.GetPublisher(pFilter)
		if err != nil {
			internalServerErrorWithSentry(c, "GetPublisher: ", err)
		}
		book.Publisher = publisher
		if publisher == nil {
			createPublisher := domain.Publisher{}
			createPublisher.Name = *form.PublisherName
			newPublisher, err := b.PublisherUseCase.CreatePublisher(createPublisher)
			if err != nil {
				log.Println("CreatePublisher: ", err.Error())
				internalServerErrorWithSentry(c, "CreatePublisher: ", err)
				return
			}
			book.Publisher = newPublisher
		}
	}
	book.Isbn = form.Isbn
	book.SmallImageUrl = form.SmallImageUrl
	book.MediumImageUrl = form.MediumImageUrl
	book.ItemUrl = form.ItemUrl
	book.AffiliateUrl = form.AffiliateUrl
	book.ReadState = domain.NotReadValue

	newBook, err := b.UseCase.CreateBook(book)
	if err != nil {
		log.Println(err.Error())
		if gin.Mode() == gin.ReleaseMode {
			sentryLogError("CreateBook: ", err)
		}
		c.JSON(http.StatusBadRequest, Response{Content: book})
		return
	}
	c.JSON(http.StatusOK, Response{Content: newBook})
}

func (b *bookController) ChangeBookStatus(c *gin.Context) {
	bookId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("ChangeBookStatus: ", err.Error())
		badRequestWithSentry(c, "ChangeBookStatus: ", err)
		return
	}
	accountId, ok := c.MustGet("account_id").(string)
	if !ok {
		log.Println("ChangeBookStatus: ")
		badRequestWithSentry(c, "ChangeBookStatus: ", errors.New("accountId parser error"))
		return
	}
	filter := usecases.NewFilter()
	usecases.ById(filter, bookId)
	usecases.ByAccountId(filter, accountId)

	book, err := b.UseCase.ChangeStatus(filter)
	if err != nil {
		log.Println(err.Error())
		internalServerErrorWithSentry(c, "ChangeBookStatus: ", err)
		return
	}
	c.JSON(http.StatusOK, Response{Content: book})
}

func (b *bookController) DeleteBook(c *gin.Context) {
	bookId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("DeleteBook: ", err.Error())
		badRequestWithSentry(c, "DeleteBook: ", err)
		return
	}
	accountId, ok := c.MustGet("account_id").(string)
	if !ok {
		log.Println("DeleteBook: ", errors.New("accountId parser error"))
		badRequestWithSentry(c, "DeleteBook: ", err)
		return
	}
	bookFilter := usecases.NewFilter()
	usecases.ById(bookFilter, bookId)
	usecases.ByAccountId(bookFilter, accountId)

	err = b.UseCase.DeleteBook(bookFilter)
	if err != nil {
		log.Println(err.Error())
		internalServerErrorWithSentry(c, "DeleteBook: ", err)
		return
	}
	c.Status(http.StatusOK)
}

func (b *bookController) UpdateBook(c *gin.Context) {
	form := BookUpdateForm{}
	err := c.ShouldBind(&form)
	if err != nil {
		log.Println("UpdateBook: ", err.Error())
		badRequestWithSentry(c, "UpdateBook: ", err)
		return
	}
	accountId, ok := c.MustGet("account_id").(string)
	if !ok {
		log.Println("UpdateBook: ", errors.New("accountId parser error"))
		badRequestWithSentry(c, "UpdateBook: ", err)
		return
	}

	filter := usecases.NewFilter()
	usecases.ByAccountId(filter, accountId)
	usecases.ById(filter, form.ID)
	book, err := b.UseCase.GetBook(filter)
	if err != nil {
		log.Println("UpdateBook: ", err.Error())
		badRequestWithSentry(c, "UpdateBook: ", err)
		return
	}

	book.Title = form.Title
	author := domain.Author{}
	if form.AuthorID != nil {
		author.ID = *form.AuthorID
		book.Author = &author
	} else {
		book.Author = nil
	}

	publisher := domain.Publisher{}
	if form.PublisherID != nil {
		publisher.ID = *form.PublisherID
		book.Publisher = &publisher
	} else {
		book.Publisher = nil
	}

	if form.StartAt != nil {
		book.StartAt = domain.NullTime{NullTime: mysql.NullTime{Time: *form.StartAt, Valid: true}}
	} else {
		book.StartAt = domain.NullTime{NullTime: mysql.NullTime{Time: domain.JstNow(), Valid: false}}
	}

	if form.EndAt != nil {
		book.EndAt = domain.NullTime{NullTime: mysql.NullTime{Time: *form.EndAt, Valid: true}}
	} else {
		book.EndAt = domain.NullTime{NullTime: mysql.NullTime{Time: domain.JstNow(), Valid: false}}
	}

	if form.StartAt == nil && form.EndAt == nil {
		book.ReadState = domain.NotReadValue
	} else if form.StartAt != nil && form.EndAt == nil {
		book.ReadState = domain.ReadingValue
	} else if form.StartAt != nil && form.EndAt == nil {
		book.ReadState = domain.ReadValue
	} else {
		internalServerErrorWithSentry(c, "UpdateBook: ", errors.New("bad read state"))
		return
	}

	updatedBook, err := b.UseCase.UpdateBook(*book, nil)
	if err != nil {
		log.Println(err.Error())
		internalServerErrorWithSentry(c, "UpdateBook: ", err)
		return
	}
	c.JSON(http.StatusOK, Response{Content: updatedBook})
}

func (b *bookController) GetCountedByAuthor(c *gin.Context) {
	filter := usecases.NewFilter()
	filter, err := addAccountToFilter(c, &filter)
	if err != nil {
		log.Println("GetCountedByAuthor: ", err.Error())
		badRequestWithSentry(c, "GetCountedByAuthor: ", err)
		return
	}

	authorCountedByName, err := b.UseCase.CountByName(filter, "author")
	if err != nil {
		log.Println("GetCountedByAuthor: ", err.Error())
		internalServerErrorWithSentry(c, "GetCountedByAuthor: ", err)
		return
	}
	c.JSON(http.StatusOK, Response{Content: authorCountedByName})
}

func (b *bookController) GetCountedByPublisher(c *gin.Context) {
	filter := usecases.NewFilter()
	filter, err := addAccountToFilter(c, &filter)
	if err != nil {
		log.Println("GetCountedByPublisher: ", err.Error())
		badRequestWithSentry(c, "GetCountedByPublisher: ", err)
		return
	}

	countedByName, err := b.UseCase.CountByName(filter, "publisher")
	if err != nil {
		log.Println("GetAllBooks: ", err.Error())
		internalServerErrorWithSentry(c, "GetCountedByPublisher: ", err)
		return
	}
	c.JSON(http.StatusOK, Response{Content: countedByName})
}

func (b *bookController) GetCountedRegisterDaily(c *gin.Context) {
	countedDate, abort := getCountedDaily(c, b, domain.BookRegister)
	if abort {
		return
	}
	c.JSON(http.StatusOK, Response{Content: countedDate})
}

func (b *bookController) GetCountedStartDaily(c *gin.Context) {
	countedDate, abort := getCountedDaily(c, b, domain.BookReadStart)
	if abort {
		return
	}
	c.JSON(http.StatusOK, Response{Content: countedDate})
}

func (b *bookController) GetCountedEndDaily(c *gin.Context) {
	countedDate, abort := getCountedDaily(c, b, domain.BookReadEnd)
	if abort {
		return
	}
	c.JSON(http.StatusOK, Response{Content: countedDate})
}

// todo ひとまず一括で取得、必要ならばpath paramで月を指定できるようにする
func getCountedDaily(c *gin.Context, b *bookController, dateKey string) (countedDate *domain.CountedDates, aborted bool) {
	defer func() { aborted = c.IsAborted() }()
	filter := usecases.NewFilter()

	filter, err := addAccountToFilter(c, &filter)
	if err != nil {
		log.Println("GetCountedDaily: ", err.Error())
		sentryLogWarning("GetCountedDaily: ", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	countedDate, err = b.UseCase.CountByDate(filter, dateKey, domain.DateKeyDaily)
	if err != nil {
		log.Println("GetCountedDaily: ", err.Error())
		sentryLogError("GetCountedDaily: ", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusInternalServerError})
		return
	}
	return
}

// todo ひとまず一括で取得、必要ならばpath paramで年を指定できるようにする
func (b *bookController) GetCountedMonthly(c *gin.Context) {
	filter := usecases.NewFilter()

	filter, err := addAccountToFilter(c, &filter)
	if err != nil {
		log.Println("GetCountedMonthly: ", err.Error())
		badRequestWithSentry(c, "GetCountedMonthly: ", err)
		return
	}

	countedDate, err := b.UseCase.CountByDate(filter, "end_at", domain.DateKeyMonthly)
	if err != nil {
		log.Println("GetCountedMonthly: ", err.Error())
		internalServerErrorWithSentry(c, "GetCountedMonthly: ", err)
		return
	}
	c.JSON(http.StatusOK, Response{Content: countedDate})
}

func addAccountToFilter(c *gin.Context, f *map[string]interface{}) (map[string]interface{}, error) {
	accountId, ok := c.MustGet("account_id").(string)
	if !ok {
		return nil, errors.New("accountId parser error")
	}
	usecases.ByAccountId(*f, accountId)
	return *f, nil
}
