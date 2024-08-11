package config

import (
    "fmt"
    "log"
    "os"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "url-shortener/models"
)

func SetupDatabase() (*gorm.DB, error) {
    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")
    dbPort := os.Getenv("DB_PORT")
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        dbHost, dbUser, dbPassword, dbName, dbPort)
   
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    // Ejecutar migraciones
    if err := db.AutoMigrate(&models.User{}, &models.URL{}); err != nil {
        return nil, err
    }

    // Crear índices manualmente
    if err := createIndices(db); err != nil {
        return nil, err
    }

    return db, nil
}

func createIndices(db *gorm.DB) error {
    // Índice para short_code (ya debe existir debido a la etiqueta uniqueIndex)
    if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code)").Error; err != nil {
        log.Printf("Warning: Could not create index on urls.short_code: %v", err)
    }

    // Índice para user_id
    if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_urls_user_id ON urls(user_id)").Error; err != nil {
        log.Printf("Warning: Could not create index on urls.user_id: %v", err)
    }

    // Índice compuesto para user_id y short_code
    if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_urls_user_id_short_code ON urls(user_id, short_code)").Error; err != nil {
        log.Printf("Warning: Could not create composite index on urls(user_id, short_code): %v", err)
    }

    return nil
}