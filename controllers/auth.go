package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "url-shortener/services"
)

type AuthController struct {
    authService services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
    return &AuthController{authService: authService}
}

type RegisterInput struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

func (ac *AuthController) Register(c *gin.Context) {
    var input RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := ac.authService.Register(input.Email, input.Password); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

type LoginInput struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

func (ac *AuthController) Login(c *gin.Context) {
    var input LoginInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    token, user, err := ac.authService.Login(input.Email, input.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user": gin.H{
            "id":    user.ID,
            "email": user.Email,
        },
    })
}