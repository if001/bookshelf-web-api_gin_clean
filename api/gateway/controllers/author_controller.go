package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"bookshelf-web-api_gin_clean/api/usecases"
	"bookshelf-web-api_gin_clean/api/gateway/repositories"
	"bookshelf-web-api_gin_clean/api/domain"
)

type AuthorForm struct {
	AuthorName string `json:"author_name" binding:"required"`
}

type authorController struct {
	UseCase usecases.AuthorUseCase
}

type AuthorController interface {
	GetCountedAuthors(c *gin.Context)
	CreateAuthor(c *gin.Context)
}

func NewAuthorController(dbConnection repositories.DBConnection) AuthorController {
	repo := repositories.NewAuthorRepository(dbConnection)
	u := usecases.NewAuthorUseCase(repo)
	return &authorController{UseCase: u}
}

func (a *authorController) GetCountedAuthors(c *gin.Context) {
	authors, err := a.UseCase.GetAllAuthor()
	if err != nil {
		log.Println("GetCountedAuthors: ", err.Error())
		internalServerErrorWithSentry(c, "GetCountedAuthors: ", err)
		return
	}
	c.JSON(http.StatusOK, Response{Content: authors})
}

func (a *authorController) CreateAuthor(c *gin.Context) {
	form := AuthorForm{}
	err := c.ShouldBind(&form)
	if err != nil {
		log.Println("CreateAuthor: ", err.Error())
		badRequestWithSentry(c, "CreateAuthor: ", err)
		return
	}

	author := domain.Author{}
	author.Name = form.AuthorName
	newAuthor, err := a.UseCase.CreateAuthor(author)
	if err != nil {
		log.Println("CreateAuthor: ", err.Error())
		internalServerErrorWithSentry(c, "CreateAuthor: ", err)
		return
	}
	c.JSON(http.StatusOK, Response{Content: newAuthor})
}
