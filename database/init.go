package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"kdf_tech_job/config"
)

var db *gorm.DB

func InitDatabaseConnection() (err error) {
	connectionsStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		config.Config.UserDatabase.Host,
		config.Config.UserDatabase.Port,
		config.Config.UserDatabase.Username,
		config.Config.UserDatabase.Password,
		config.Config.UserDatabase.Name,
	)
	db, err = gorm.Open(postgres.Open(connectionsStr), &gorm.Config{})
	return nil
}
