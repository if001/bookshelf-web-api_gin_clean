package database

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DB   DBConf
	Addr string `envconfig:"port" default:":8080"`
}

type DBConf struct {
	User     string `envconfig:"mysql_user" default:"api"`
	Password string `envconfig:"mysql_password" default:"hogehoge"`
	Host     string `envconfig:"mysql_ip" default:"127.0.0.1:3306"`
	DB       string `envconfig:"mysql_db" default:"bookshelf"`
}

func LoadConfig() (*Config, error) {
	var config = Config{}
	if err := envconfig.Process("APP", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
