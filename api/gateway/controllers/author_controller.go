package controllers

import (
	"bookshelf-web-api_gin_clean/api/usecases"
	"github.com/gin-gonic/gin"
	"bookshelf-web-api_gin_clean/api/gateway/repositories"
	"log"
	"net/http"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, Response{Content: authors})
}

func (a *authorController) CreateAuthor(c *gin.Context) {
	form := AuthorForm{}
	err := c.ShouldBind(&form)
	if err != nil {
		log.Println("CreateAuthor: ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	author := domain.Author{}
	author.Name = form.AuthorName
	newAuthor, err := a.UseCase.CreateAuthor(author)
	if err != nil {
		log.Println("CreateAuthor: ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}
	c.JSON(http.StatusOK, Response{Content: newAuthor})
}
