package externalInteface

import (
	"bookshelf-web-api_gin_clean/api/externalInteface/database"
	"bookshelf-web-api_gin_clean/api/gateway/controllers"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()
	router.Use(Options, authMiddleware())
	//router.Use(Options, authMiddlewareTest())
	conn := database.NewSqlConnection()

	b := controllers.NewBookController(&conn)
	d := controllers.NewDescriptionController(&conn)
	a := controllers.NewAuthorController(&conn)

	router.GET("/books", b.GetAllBooks)
	router.POST("/books", b.CreateBook)
	router.PUT("/books", b.UpdateBook)

	router.GET("/book/:id", b.GetBook)
	router.DELETE("/book/:id", b.DeleteBook)

	router.PUT("/book/:id/state/start", b.ChangeBookStatus)
	router.PUT("/book/:id/state/end", b.ChangeBookStatus)

	router.GET("/book/:id/description", d.GetAllDescriptions)
	router.POST("/book/:id/description", d.CreateDescription)
	router.DELETE("/description/:id", d.DeleteDescription)

	router.GET("/counted_authors", a.GetCountedAuthors)
	router.POST("/author", a.CreateAuthor)

	return router
}
