package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "url-shortener/services"
)

type URLController struct {
    urlService services.URLService
}

func NewURLController(urlService services.URLService) *URLController {
    return &URLController{urlService: urlService}
}

func (uc *URLController) ShortenURL(c *gin.Context) {
    var input struct {
        LongURL  string `json:"long_url" binding:"required,url"`
        Campaign string `json:"campaign"`
        Medium   string `json:"medium"`
        Source   string `json:"source"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
    url, err := uc.urlService.ShortenURL(input.LongURL, input.Campaign, input.Medium, input.Source, userID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create shortened URL"})
        return
    }
    c.JSON(http.StatusCreated, gin.H{
        "message":    "URL shortened successfully",
        "short_code": url.ShortCode,
        "campaign":   url.Campaign,
        "medium":     url.Medium,
        "source":     url.Source,
    })
}

func (uc *URLController) RedirectURL(c *gin.Context) {
    shortCode := c.Param("shortCode")
   
    url, err := uc.urlService.GetURLByShortCode(shortCode)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
        return
    }
    err = uc.urlService.IncrementViewCount(url)
    if err != nil {
        // Log the error, but continue with redirection
        c.Error(err)
    }
    c.Redirect(http.StatusFound, url.LongURL)
}

func (uc *URLController) GetURLStats(c *gin.Context) {
    shortCode := c.Param("shortCode")
    userID, _ := c.Get("user_id")
    url, err := uc.urlService.GetURLStats(shortCode, userID.(uint))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found or you don't have permission to view it"})
        return
    }
    c.JSON(http.StatusOK, url)
}

func (uc *URLController) GetUserURLs(c *gin.Context) {
    userID, _ := c.Get("user_id")
    urls, err := uc.urlService.GetUserURLs(userID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user's URLs"})
        return
    }
    c.JSON(http.StatusOK, urls)
}

func (uc *URLController) GetURLsByCampaign(c *gin.Context) {
    campaign := c.Query("campaign")
    userID, _ := c.Get("user_id")
    urls, err := uc.urlService.GetURLsByCampaign(campaign, userID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch URLs"})
        return
    }
    c.JSON(http.StatusOK, urls)
}