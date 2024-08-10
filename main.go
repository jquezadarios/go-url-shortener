package main

import (
    "log"
    "os"
    "github.com/gin-gonic/gin"
    "url-shortener/config"
    "url-shortener/controllers"
    "url-shortener/middleware"
    "url-shortener/models"
)

func main() {
    db, err := config.SetupDatabase()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Migrate the schema
    err = db.AutoMigrate(&models.User{})
    if err != nil {
        log.Fatalf("Failed to migrate database schema: %v", err)
    }

    // Initialize JWT key
    middleware.InitJWTKey(os.Getenv("JWT_SECRET_KEY"))

    r := gin.Default()

    authController := controllers.NewAuthController(db)

    auth := r.Group("/auth")
    {
        auth.POST("/register", authController.Register)
        auth.POST("/login", authController.Login)
    }

    api := r.Group("/api")
    api.Use(middleware.JWTMiddleware())
    {
        // Protected routes here
    }

    r.Run(":8080")
}