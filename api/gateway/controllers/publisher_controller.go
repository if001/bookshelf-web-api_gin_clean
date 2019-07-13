package controllers

import (
	"bookshelf-web-api_gin_clean/api/domain"
	"bookshelf-web-api_gin_clean/api/gateway/repositories"
	"bookshelf-web-api_gin_clean/api/usecases"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type PublisherForm struct {
	// PublisherName string `json:"publisher_name" binding:"required"`
	PublisherName string `json:"publisher_name"`
}

type publisherController struct {
	UseCase usecases.PublisherUseCase
}

type PublisherController interface {
	GetCountedPublishers(c *gin.Context)
	CreatePublisher(c *gin.Context)
}

func NewPublisherController(dbConnection repositories.DBConnection) PublisherController {
	repo := repositories.NewPublisherRepository(dbConnection)
	u := usecases.NewPublisherUseCase(repo)
	return &publisherController{UseCase: u}
}

func (p *publisherController) GetCountedPublishers(c *gin.Context) {
	publishers, err := p.UseCase.GetAllPublisher()
	if err != nil {
		log.Println("GetCountedPublishers: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, Response{Content: publishers})
}

func (p *publisherController) CreatePublisher(c *gin.Context) {
	form := PublisherForm{}
	err := c.ShouldBind(&form)
	if err != nil {
		log.Println("CreatePublisher: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	publisher := domain.Publisher{}
	publisher.Name = form.PublisherName
	newPublisher, err := p.UseCase.CreatePublisher(publisher)
	if err != nil {
		log.Println("CreatePublisher: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusOK, Response{Content: newPublisher})
}
