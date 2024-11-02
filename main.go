package main

import (
	"noteapp/controllers"
	"noteapp/models"
	"noteapp/routes"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDatabase() *gorm.DB {
	// Параметры подключения к базе данных
	dsn := "host=localhost user=postgres password=root dbname=noteapp port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Миграция моделей
	db.AutoMigrate(&models.Note{})
	return db
}

func main() {
	db := InitDatabase()
	controllers.SetDatabase(db) // Устанавливаем базу данных в контроллере

	r := gin.Default()
	routes.RegisterRoutes(r)

	// Запуск сервера
	r.Run(":8080")
}
