package controllers

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"fmt"
)

func GetPaginate(c *gin.Context) (uint64, uint64, error){
	var page uint64 = 0
	pageStr := c.Query("page")
	if pageStr != "" {
		tmpPage, err := strconv.ParseUint(pageStr, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("GetPaginate: %s",err)
		}
		page = tmpPage
	}

	var perPage uint64 = 0
	perPageStr := c.Query("per_page")
	if perPageStr != "" {
		tmpPerPage, err := strconv.ParseUint(perPageStr, 10, 64)
		if err != nil {
			return 0, 0, fmt.Errorf("GetPaginate: %s",err)
		}
		perPage = tmpPerPage
	}
	return page, perPage, nil
}