package config

import (
	"encoding/json"
	"os"
)

type UserDatabase struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type ApplicationConfig struct {
	ListenPort   string       `json:"port"`
	UserDatabase UserDatabase `json:"database"`
}

var Config *ApplicationConfig

func LoadConfiguration(fileName string) error {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return err
	}
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config = &ApplicationConfig{}
	err = decoder.Decode(Config)
	if err != nil {
		return err
	}
	return nil
}
