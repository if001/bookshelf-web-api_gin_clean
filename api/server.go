package main

import (
	"bookshelf-web-api_gin_clean/api/externalInteface"
	"fmt"
	"os"
)

func main() {
	var port string

	port = os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	err := externalInteface.Router().Run(fmt.Sprintf(":%v",port))
	if err != nil {
		panic(fmt.Errorf("[FAILED] start sever. err: %v", err))
	}
}
