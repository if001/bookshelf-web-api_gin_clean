package database

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DB   DBConf
	Addr string `envconfig:"port" default:":8080"`
}

type DBConf struct {
	User     string `envconfig:"MYSQL_USER" default:"api"`
	Password string `envconfig:"MYSQL_PASSWORD" default:"hogehoge"`
	Host     string `envconfig:"MYSQL_IP" default:"127.0.0.1:3306"`
	DB       string `envconfig:"MYSQL_DB" default:"bookshelf"`
}

func LoadConfig() (*Config, error) {
	var config = Config{}
	if err := envconfig.Process("APP", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
