package main

import (
    "log"
    "os"
    "github.com/gin-gonic/gin"
    "github.com/rs/cors"
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

    // Configurar CORS
    config := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000"}, // Ajusta esto a la URL de tu frontend
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
        AllowCredentials: true,
    })

    // Usar el middleware CORS
    r.Use(func(c *gin.Context) {
        config.HandlerFunc(c.Writer, c.Request)
        c.Next()
    })

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