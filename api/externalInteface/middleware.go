package externalInteface

import (
	"context"
	"firebase.google.com/go"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"os"
	"strings"
)

func authMiddlewareTest() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("account_id", "UbazgLQ5MBafKlZukidXLC3a97f1")
		c.Next()
	}
}
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Firebase SDK のセットアップ
		ctx := context.Background()
		authKey := os.Getenv("FIREBASE_KEYFILE_JSON")
		var opt option.ClientOption
		if authKey == "" {
			opt = option.WithCredentialsFile("/Users/issei/gcloud_key_json/bookshelf-239408-firebase-adminsdk-ujfj8-61a6ff4292.json")
		} else {
			credentials, err := google.CredentialsFromJSON(ctx, []byte(authKey))
			if err != nil {
				fmt.Printf("authMiddleware: %v\n", err)
				os.Exit(1)
			}
			opt = option.WithCredentials(credentials)
		}

		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			fmt.Printf("authMiddleware: %v\n", err)
			os.Exit(1)
		}
		auth, err := app.Auth(context.Background())
		if err != nil {
			fmt.Printf("authMiddleware: %v\n", err)
			os.Exit(1)
		}

		// クライアントから送られてきた JWT 取得
		authHeader := c.GetHeader(("Authorization"))
		idToken := strings.Replace(authHeader, "Bearer ", "", 1)

		// JWT の検証
		token, err := auth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			// JWT が無効なら Handler に進まず別処理
			fmt.Printf("error verifying ID token: %v\n", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		log.Printf("Verified ID token: %v\n", token)
		c.Set("account_id", token.UID)
		c.Next()
	}
}

func Options() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		if c.Request.Method == "OPTIONS" {
			fmt.Println("option")
			c.AbortWithStatus(http.StatusOK)
			return
		} else {
			fmt.Println("not option")
			c.Next()
		}
	}
}