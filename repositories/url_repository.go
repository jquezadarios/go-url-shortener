package repositories

import (
    "url-shortener/models"
    "gorm.io/gorm"
)

type URLRepository interface {
    Create(url *models.URL) error
    FindByShortCode(shortCode string) (*models.URL, error)
    IncrementViewCount(url *models.URL) error
    FindByShortCodeAndUserID(shortCode string, userID uint) (*models.URL, error)
    FindByUserID(userID uint) ([]models.URL, error)
    FindByCampaignAndUserID(campaign string, userID uint) ([]models.URL, error)
    FindByMediumAndUserID(medium string, userID uint) ([]models.URL, error)
    FindBySourceAndUserID(source string, userID uint) ([]models.URL, error)
}

type urlRepository struct {
    db *gorm.DB
}

func NewURLRepository(db *gorm.DB) URLRepository {
    return &urlRepository{db: db}
}

func (r *urlRepository) Create(url *models.URL) error {
    return r.db.Create(url).Error
}

func (r *urlRepository) FindByShortCode(shortCode string) (*models.URL, error) {
    var url models.URL
    err := r.db.Where("short_code = ?", shortCode).First(&url).Error
    return &url, err
}

func (r *urlRepository) IncrementViewCount(url *models.URL) error {
    return r.db.Model(url).Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func (r *urlRepository) FindByShortCodeAndUserID(shortCode string, userID uint) (*models.URL, error) {
    var url models.URL
    err := r.db.Where("short_code = ? AND user_id = ?", shortCode, userID).First(&url).Error
    return &url, err
}

func (r *urlRepository) FindByUserID(userID uint) ([]models.URL, error) {
    var urls []models.URL
    err := r.db.Where("user_id = ?", userID).Find(&urls).Error
    return urls, err
}

func (r *urlRepository) FindByCampaignAndUserID(campaign string, userID uint) ([]models.URL, error) {
    var urls []models.URL
    err := r.db.Where("campaign = ? AND user_id = ?", campaign, userID).Find(&urls).Error
    return urls, err
}

func (r *urlRepository) FindByMediumAndUserID(medium string, userID uint) ([]models.URL, error) {
    var urls []models.URL
    err := r.db.Where("medium = ? AND user_id = ?", medium, userID).Find(&urls).Error
    return urls, err
}

func (r *urlRepository) FindBySourceAndUserID(source string, userID uint) ([]models.URL, error) {
    var urls []models.URL
    err := r.db.Where("source = ? AND user_id = ?", source, userID).Find(&urls).Error
    return urls, err
}