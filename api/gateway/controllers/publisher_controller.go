package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"bookshelf-web-api_gin_clean/api/domain"
	"bookshelf-web-api_gin_clean/api/gateway/repositories"
	"bookshelf-web-api_gin_clean/api/usecases"
)

type PublisherForm struct {
	PublisherName string `json:"publisher_name" binding:"required"`
	// PublisherName string `json:"publisher_name"`
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
		internalServerErrorWithSentry(c, "GetCountedPublishers: ", err)
		return
	}
	c.JSON(http.StatusOK, Response{Content: publishers})
}

func (p *publisherController) CreatePublisher(c *gin.Context) {
	form := PublisherForm{}
	err := c.ShouldBind(&form)
	if err != nil {
		log.Println("CreatePublisher: ", err.Error())
		badRequestWithSentry(c, "CreatePublisher: ", err)
		return
	}

	publisher := domain.Publisher{}
	publisher.Name = form.PublisherName
	newPublisher, err := p.UseCase.CreatePublisher(publisher)
	if err != nil {
		log.Println("CreatePublisher: ", err.Error())
		internalServerErrorWithSentry(c, "CreatePublisher: ", err)
		return
	}
	c.JSON(http.StatusOK, Response{Content: newPublisher})
}
