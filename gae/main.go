package main

import (
	"bookshelf-web-api_gin_clean/api/externalInteface"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"net/http"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := externalInteface.Router()
	http.Handle("/", r)
	appengine.Main()
}
