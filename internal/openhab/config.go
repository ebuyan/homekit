package openhab

import (
	"encoding/base64"
	"os"
)

type Config struct {
	login     string
	password  string
	searchTag string
}

func NewConfig() Config {
	return Config{
		login:     os.Getenv("OPENHAB_LOGIN"),
		password:  os.Getenv("OPENHAB_PWD"),
		searchTag: "homekit",
	}
}

func (c Config) GetCredentials() string {
	auth := c.login + ":" + c.password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
