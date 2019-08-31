package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
)

var tagInternalServerError = map[string]string{
	"tag":         "bookshelf-web-api_gin_clean",
	"status_code": string(http.StatusInternalServerError),
}

var tagStatusBadRequest = map[string]string{
	"tag":         "bookshelf-web-api_gin_clean",
	"status_code": string(http.StatusBadRequest),
}

func sentryLogError(funcName string, err error) {
	msg := errors.New(fmt.Sprintf("%s: %s", funcName, err.Error()))
	raven.CaptureError(msg, tagInternalServerError)
}

func sentryLogWarning(funcName string, err error) {
	p := raven.NewPacket(fmt.Sprintf("%s: %s", funcName, err.Error()))
	p.Level = raven.WARNING
	raven.Capture(p, tagStatusBadRequest)
}

func badRequestWithSentry(c *gin.Context, funcName string, err error) {
	if gin.Mode() == gin.ReleaseMode {
		sentryLogWarning(funcName, err)
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
}

func internalServerErrorWithSentry(c *gin.Context, funcName string, err error) {
	if gin.Mode() == gin.ReleaseMode {
		sentryLogError(funcName, err)
	}
	c.JSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
}
