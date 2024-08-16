package services

import (
    "url-shortener/models"
    "url-shortener/repositories"
    "url-shortener/cache"
    "github.com/teris-io/shortid"
    "fmt"
)

type URLService interface {
    ShortenURL(longURL string, name string,campaign string, medium string, source string, userID uint) (*models.URL, error)
    GetURLByShortCode(shortCode string) (*models.URL, error)
    IncrementViewCount(url *models.URL) error
    GetURLStats(shortCode string, userID uint) (*models.URL, error)
    GetUserURLs(userID uint) ([]models.URL, error)
    GetURLsByCampaign(campaign string, userID uint) ([]models.URL, error)
    GetURLsByMedium(medium string, userID uint) ([]models.URL, error)
    GetURLsBySource(source string, userID uint) ([]models.URL, error)
}

type urlService struct {
    urlRepo repositories.URLRepository
    cache   *cache.MemcachedClient
}

func NewURLService(urlRepo repositories.URLRepository, cache *cache.MemcachedClient) URLService {
    return &urlService{urlRepo: urlRepo, cache: cache}
}

func (s *urlService) ShortenURL(longURL string, name string, campaign string, medium string, source string, userID uint) (*models.URL, error) {
    shortCode, err := shortid.Generate()
    if err != nil {
        return nil, err
    }
    url := &models.URL{
        LongURL:   longURL,
        ShortCode: shortCode,
        UserID:    userID,
        Name:      name,
        Campaign:  campaign,
        Medium:    medium,
        Source:    source,
    }
    err = s.urlRepo.Create(url)
    if err != nil {
        return nil, err
    }
    // Cache the new URL
    cacheKey := fmt.Sprintf("url:%s", shortCode)
    s.cache.Set(cacheKey, url, 3600) // Cache for 1 hour
    return url, nil
}

func (s *urlService) GetURLByShortCode(shortCode string) (*models.URL, error) {
    cacheKey := fmt.Sprintf("url:%s", shortCode)
   
    // Try to get from cache first
    var url *models.URL
    err := s.cache.Get(cacheKey, &url)
    if err == nil {
        return url, nil
    }
    // If not in cache, get from database
    url, err = s.urlRepo.FindByShortCode(shortCode)
    if err != nil {
        return nil, err
    }
    // Cache the result
    s.cache.Set(cacheKey, url, 3600) // Cache for 1 hour
    return url, nil
}

func (s *urlService) IncrementViewCount(url *models.URL) error {
    err := s.urlRepo.IncrementViewCount(url)
    if err != nil {
        return err
    }
    // Update cache
    cacheKey := fmt.Sprintf("url:%s", url.ShortCode)
    s.cache.Delete(cacheKey) // Delete old cache
    s.cache.Set(cacheKey, url, 3600) // Set new cache
    return nil
}

func (s *urlService) GetURLStats(shortCode string, userID uint) (*models.URL, error) {
    cacheKey := fmt.Sprintf("stats:%s:%d", shortCode, userID)
    var url *models.URL
    err := s.cache.Get(cacheKey, &url)
    if err == nil {
        return url, nil
    }
    url, err = s.urlRepo.FindByShortCodeAndUserID(shortCode, userID)
    if err != nil {
        return nil, err
    }
    s.cache.Set(cacheKey, url, 3600) // Cache for 1 hour
    return url, nil
}

func (s *urlService) GetUserURLs(userID uint) ([]models.URL, error) {
    cacheKey := fmt.Sprintf("user_urls:%d", userID)
    var urls []models.URL
    err := s.cache.Get(cacheKey, &urls)
    if err == nil {
        return urls, nil
    }
    urls, err = s.urlRepo.FindByUserID(userID)
    if err != nil {
        return nil, err
    }
    s.cache.Set(cacheKey, urls, 3600) // Cache for 1 hour
    return urls, nil
}

func (s *urlService) GetURLsByCampaign(campaign string, userID uint) ([]models.URL, error) {
    cacheKey := fmt.Sprintf("campaign_urls:%s:%d", campaign, userID)
    var urls []models.URL
    err := s.cache.Get(cacheKey, &urls)
    if err == nil {
        return urls, nil
    }
    urls, err = s.urlRepo.FindByCampaignAndUserID(campaign, userID)
    if err != nil {
        return nil, err
    }
    s.cache.Set(cacheKey, urls, 3600) // Cache for 1 hour
    return urls, nil
}
func (s *urlService) GetURLsByMedium(medium string, userID uint) ([]models.URL, error) {
    cacheKey := fmt.Sprintf("medium_urls:%s:%d", medium, userID)
    var urls []models.URL
    err := s.cache.Get(cacheKey, &urls)
    if err == nil {
        return urls, nil
    }
    urls, err = s.urlRepo.FindByMediumAndUserID(medium, userID)
    if err != nil {
        return nil, err
    }
    s.cache.Set(cacheKey, urls, 3600) // Cache for 1 hour
    return urls, nil
}

func (s *urlService) GetURLsBySource(source string, userID uint) ([]models.URL, error) {
    cacheKey := fmt.Sprintf("source_urls:%s:%d", source, userID)
    var urls []models.URL
    err := s.cache.Get(cacheKey, &urls)
    if err == nil {
        return urls, nil
    }
    urls, err = s.urlRepo.FindBySourceAndUserID(source, userID)
    if err != nil {
        return nil, err
    }
    s.cache.Set(cacheKey, urls, 3600) // Cache for 1 hour
    return urls, nil
}