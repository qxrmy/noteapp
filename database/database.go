package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"noteapp/config"
	"noteapp/models"
)

var DB *gorm.DB

func InitDatabase() {
	cfg := config.LoadConfig()

	dsn := "host=" + cfg.DBHost +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" port=" + cfg.DBPort +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}

	db.AutoMigrate(&models.Note{})

	DB = db
}
