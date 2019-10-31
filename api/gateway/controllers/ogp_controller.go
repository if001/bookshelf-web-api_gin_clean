package controllers

import (
	"bookshelf-web-api_gin_clean/api/gateway/repositories"
	"bookshelf-web-api_gin_clean/api/usecases"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewOgpController(dbConnection repositories.DBConnection) OgpController {
	repo := repositories.NewBookRepository(dbConnection)
	u := usecases.NewBookUseCase(repo)
	return &ogpController{UseCase: u}
}

type ogpController struct {
	UseCase usecases.BookUseCase
}

type OgpController interface {
	TemplateWithOGPHeader(c *gin.Context)
}

const shareURL = "https://bookstorage.edgwbs.net/share/%d"

func (o *ogpController) TemplateWithOGPHeader(c *gin.Context) {
	bookId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		log.Println("TemplateWithOGPHeader: ", err.Error())
		badRequestWithSentry(c, "TemplateWithOGPHeader: ", err)
		return
	}

	ua := c.GetHeader("user-agent")
	if isTwitter(ua) {
		filter := usecases.NewFilter()
		usecases.ById(filter, bookId)
		book, err := o.UseCase.GetBook(filter)
		if err != nil {
			internalServerErrorWithSentry(c, "TemplateWithOGPHeader: ", err)
			return
		}
		obj := map[string]interface{}{
			"SiteURL":  fmt.Sprintf(shareURL, book.ID),
			"ImageURL": book.MediumImageUrl,
		}
		c.Header("Content-Type", "text/html")
		c.HTML(http.StatusOK, "ogp_header.html", obj)
		c.Abort()
		return
	} else {
		fmt.Println(fmt.Sprintf(shareURL, bookId))
		c.Redirect(http.StatusMovedPermanently, fmt.Sprintf(shareURL, bookId))
		c.Abort()
		return
	}
}

func isTwitter(ua string) bool {
	return ua == "twitterbot/1.0" || ua == "Twitterbot" || ua == "twitterbot" || ua == "Twitterbot/1.0"
}
