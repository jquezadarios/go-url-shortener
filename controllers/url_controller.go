package controllers

import (
	"net/http"
	"url-shortener/models"
	"log"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
	"gorm.io/gorm"
)

type URLController struct {
	db *gorm.DB
}

func NewURLController(db *gorm.DB) *URLController {
	return &URLController{db: db}
}

func (uc *URLController) ShortenURL(c *gin.Context) {
	var input struct {
		LongURL string `json:"long_url" binding:"required,url"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from the JWT token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Generate short code
	shortCode, err := shortid.Generate()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate short code"})
		return
	}

	url := models.URL{
		LongURL:   input.LongURL,
		ShortCode: shortCode,
		UserID:    userID.(uint),
	}

	if err := uc.db.Create(&url).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shortened URL"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "URL shortened successfully",
		"short_code": url.ShortCode,
	})
}

func (uc *URLController) RedirectURL(c *gin.Context) {
    shortCode := c.Param("shortCode")

    var url models.URL
    if err := uc.db.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
        return
    }

    // Incrementar el contador de vistas en una transacción
    err := uc.db.Transaction(func(tx *gorm.DB) error {
        result := tx.Model(&url).Update("view_count", gorm.Expr("view_count + ?", 1))
        if result.Error != nil {
            return result.Error
        }
        if result.RowsAffected == 0 {
            return fmt.Errorf("no rows affected when updating view count")
        }
        return nil
    })

    if err != nil {
        log.Printf("Failed to update view count: %v", err)
    } else {
        log.Printf("Updated view count for short code %s", shortCode)
    }

    // Obtener la URL actualizada
    if err := uc.db.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
        log.Printf("Failed to fetch updated URL: %v", err)
    } else {
        log.Printf("Current view count for short code %s: %d", shortCode, url.ViewCount)
    }

    // Realizar la redirección
    c.Redirect(http.StatusFound, url.LongURL)
}

func (uc *URLController) GetURLStats(c *gin.Context) {
	shortCode := c.Param("shortCode")
	userID, _ := c.Get("user_id")

	var url models.URL
	if err := uc.db.Where("short_code = ? AND user_id = ?", shortCode, userID).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found or you don't have permission to view it"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_code": url.ShortCode,
		"long_url":   url.LongURL,
		"view_count": url.ViewCount,
		"created_at": url.CreatedAt,
	})
}

func (uc *URLController) GetUserURLs(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var urls []models.URL
	if err := uc.db.Where("user_id = ?", userID).Find(&urls).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user's URLs"})
		return
	}

	c.JSON(http.StatusOK, urls)
}

