package externalInteface

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB DBConf
}

type DBConf struct {
	User     string
	Password string
	Host     string
	DB       string
}

func (d DBConf) getURL() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		d.User,
		d.Password,
		d.Host,
		d.DB)
}

func LoadConfig() Config {
	var config = Config{}
	goEnv := os.Getenv("GO_ENV")
	if goEnv == "" {
		goEnv = "local"
	}
	envFile := fmt.Sprintf(".env.%s", goEnv)
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file", envFile)
	}
	config.DB.User = os.Getenv("APP_MYSQL_USER")
	config.DB.Password = os.Getenv("APP_MYSQL_PASSWORD")
	config.DB.Host = os.Getenv("APP_MYSQL_IP")
	config.DB.DB = os.Getenv("APP_MYSQL_DB")

	return config
}
