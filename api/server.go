package main

import (
	"bookshelf-web-api_gin_clean/api/externalInteface"
	"fmt"
	"os"
)

func main() {
	var addr string

	addr = os.Getenv("PORT")
	if addr == "" {
		addr = "8081"
	}

	err := externalInteface.Router().Run(addr)
	if err != nil {
		panic(fmt.Errorf("[FAILED] start sever. err: %v", err))
	}
}
