package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"noteapp/config" // Импортируйте пакет config
	"noteapp/models"
)

var DB *gorm.DB

func InitDatabase() {
	cfg := config.LoadConfig() // Загружаем конфигурацию

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

	// Миграции
	db.AutoMigrate(&models.Note{})

	DB = db
}
