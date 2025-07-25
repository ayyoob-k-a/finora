package db

import (
	"fmt"

	"github.com/ayyoob-k-a/finora/configs"
	"github.com/ayyoob-k-a/finora/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg configs.Config)  (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s", cfg.DBUSER, cfg.DBNAME, cfg.PASSWORD, cfg.HOST, cfg.PORT)
	// psqlInfo:="user=postgres dbname=company password=123 host=localhost port=5432 "
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})
	if err != nil {
		fmt.Println("error in db connection", err)
	}

	db.AutoMigrate(&domain.User{})

	// Create a admin in postgress database using terminal with credential of name and password;
	// db.AutoMigrate(&domain.Admin{})


	return db, nil

}
