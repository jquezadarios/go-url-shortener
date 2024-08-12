package services

import (
    "url-shortener/models"
    "url-shortener/repositories"
    "github.com/teris-io/shortid"
)

type URLService interface {
    ShortenURL(longURL string, campaign string, userID uint) (*models.URL, error)
    GetURLByShortCode(shortCode string) (*models.URL, error)
    IncrementViewCount(url *models.URL) error
    GetURLStats(shortCode string, userID uint) (*models.URL, error)
    GetUserURLs(userID uint) ([]models.URL, error)
    GetURLsByCampaign(campaign string, userID uint) ([]models.URL, error)
}

type urlService struct {
    urlRepo repositories.URLRepository
}

func NewURLService(urlRepo repositories.URLRepository) URLService {
    return &urlService{urlRepo: urlRepo}
}

func (s *urlService) ShortenURL(longURL string, campaign string, userID uint) (*models.URL, error) {
    shortCode, err := shortid.Generate()
    if err != nil {
        return nil, err
    }

    url := &models.URL{
        LongURL:   longURL,
        ShortCode: shortCode,
        UserID:    userID,
        Campaign:  campaign,
    }

    err = s.urlRepo.Create(url)
    if err != nil {
        return nil, err
    }

    return url, nil
}

func (s *urlService) GetURLByShortCode(shortCode string) (*models.URL, error) {
    return s.urlRepo.FindByShortCode(shortCode)
}

func (s *urlService) IncrementViewCount(url *models.URL) error {
    return s.urlRepo.IncrementViewCount(url)
}

func (s *urlService) GetURLStats(shortCode string, userID uint) (*models.URL, error) {
    return s.urlRepo.FindByShortCodeAndUserID(shortCode, userID)
}

func (s *urlService) GetUserURLs(userID uint) ([]models.URL, error) {
    return s.urlRepo.FindByUserID(userID)
}

func (s *urlService) GetURLsByCampaign(campaign string, userID uint) ([]models.URL, error) {
    return s.urlRepo.FindByCampaignAndUserID(campaign, userID)
}