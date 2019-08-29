package externalInteface

import (
	"bookshelf-web-api_gin_clean/api/externalInteface/database"
	"bookshelf-web-api_gin_clean/api/gateway/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery(), Options())

	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, "ok"); return })

	config := LoadConfig()

	conn := database.NewSqlConnection(config.DB.getURL())

	b := controllers.NewBookController(&conn)
	d := controllers.NewDescriptionController(&conn)
	a := controllers.NewAuthorController(&conn)
	p := controllers.NewPublisherController(&conn)

	authorized := router.Group("/")
	authorized.Use(authMiddleware())
	// authorized.Use(authMiddlewareTest())
	{
		authorized.GET("/books", b.GetAllBooks)
		authorized.POST("/books", b.CreateBook)
		authorized.PUT("/books", b.UpdateBook)
		authorized.GET("/books/counted/author", b.GetCountedByAuthor)
		authorized.GET("/books/counted/publisher", b.GetCountedByPublisher)
		authorized.GET("/books/counted/daily", b.GetCountedDaily)
		authorized.GET("/books/counted/monthly", b.GetCountedMonthly)

		authorized.GET("/book/:id", b.GetBook)
		authorized.DELETE("/book/:id", b.DeleteBook)

		authorized.PUT("/book/:id/state/start", b.ChangeBookStatus)
		authorized.PUT("/book/:id/state/end", b.ChangeBookStatus)

		authorized.GET("/book/:id/description", d.GetAllDescriptions)
		authorized.POST("/book/:id/description", d.CreateDescription)
		authorized.DELETE("/description/:id", d.DeleteDescription)

		authorized.GET("/counted_authors", a.GetCountedAuthors)
		authorized.POST("/author", a.CreateAuthor)

		authorized.GET("/counted_publisher", p.GetCountedPublishers)
		authorized.POST("/publisher", p.CreatePublisher)
	}
	return router
}
