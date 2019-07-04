package controllers

import (
	"github.com/gin-gonic/gin"
	"bookshelf-web-api_gin_clean/api/usecases"
	"bookshelf-web-api_gin_clean/api/gateway/repositories"
	"strconv"
	"log"
	"net/http"
	"bookshelf-web-api_gin_clean/api/domain"
)

type descriptionController struct {
	UseCase usecases.DescriptionUseCase
}

type DescriptionController interface {
	GetAllDescriptions(c *gin.Context)
	CreateDescription(c *gin.Context)
	DeleteDescription(c *gin.Context)
}

func NewDescriptionController(dbConnection repositories.DBConnection) DescriptionController {
	descRepo := repositories.NewDescriptionRepository(dbConnection)
	bookRepo := repositories.NewBookRepository(dbConnection)
	u := usecases.NewDescriptionUseCase(descRepo, bookRepo)
	return &descriptionController{UseCase: u}
}

type DescriptionForm struct {
	BookID  uint64 `json:"book_id"`
	Content string `json:"content"`
}

func (d descriptionController) GetAllDescriptions(c *gin.Context) {
	bookId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("GetAllDescriptions: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	page, perPage, err := GetPaginate(c)
	if err != nil {
		log.Println("GetPaginate: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusInternalServerError})
		return
	}

	filter := usecases.NewFilter()
	usecases.ByBookId(filter, bookId)

	description, err := d.UseCase.GetAllDescriptions(filter, page, perPage)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusNotFound)})
		return
	}
	c.JSON(http.StatusOK, Response{Content: description})
}

func (d descriptionController) CreateDescription(c *gin.Context) {
	bookId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("GetAllDescriptions: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	form := DescriptionForm{}
	err = c.ShouldBind(&form)
	if err != nil {
		log.Println("GetAllDescriptions: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	description := domain.Description{
		BookId:  bookId,
		Content: form.Content,
	}

	newDescription, err := d.UseCase.CreateDescription(description)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusNotFound)})
		return
	}
	c.JSON(http.StatusOK, Response{Content: newDescription})
}

func (d descriptionController) DeleteDescription(c *gin.Context) {
	descriptionId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("GetAllDescriptions: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}
	description := domain.Description{}
	description.ID = descriptionId

	err = d.UseCase.DeleteDescription(description)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusNotFound)})
		return
	}
	c.Status(http.StatusOK)
}
