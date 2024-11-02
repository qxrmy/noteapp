package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"

	"noteapp/models"
)

var db *gorm.DB

func SetDatabase(database *gorm.DB) {
	db = database
}

func CreateNote(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
		return
	}

	var note models.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if note.Title == "" || note.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and Content cannot be empty"})
		return
	}

	if err := db.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}
	c.JSON(http.StatusOK, note)
}

func GetNotes(c *gin.Context) {
	var notes []models.Note
	var total int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	offset := (page - 1) * pageSize

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	query := db.Model(&models.Note{})

	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count notes"})
		return
	}

	if err := query.Limit(pageSize).Offset(offset).Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve notes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"notes":     notes,
	})
}

func GetNoteByID(c *gin.Context) {
	id := c.Param("id")
	var note models.Note
	if err := db.First(&note, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}
	c.JSON(http.StatusOK, note)
}

func UpdateNote(c *gin.Context) {
	id := c.Param("id")
	var note models.Note

	if err := db.First(&note, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if note.Title == "" || note.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and Content cannot be empty"})
		return
	}

	if err := db.Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update note"})
		return
	}

	c.JSON(http.StatusOK, note)
}

func DeleteNote(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&models.Note{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}

	db.Exec("ALTER SEQUENCE notes_id_seq RESTART WITH 1")

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
}
