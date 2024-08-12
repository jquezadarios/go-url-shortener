package services

import (
    "errors"
    "url-shortener/models"
    "url-shortener/repositories"
    "url-shortener/middleware"
    "golang.org/x/crypto/bcrypt"
)

type AuthService interface {
    Register(email, password string) error
    Login(email, password string) (string, *models.User, error)
}

type authService struct {
    userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
    return &authService{userRepo: userRepo}
}

func (s *authService) Register(email, password string) error {
    if _, err := s.userRepo.FindByEmail(email); err == nil {
        return errors.New("email already exists")
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user := &models.User{
        Email:    email,
        Password: string(hashedPassword),
    }

    return s.userRepo.Create(user)
}

func (s *authService) Login(email, password string) (string, *models.User, error) {
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return "", nil, errors.New("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", nil, errors.New("invalid credentials")
    }

    token, err := middleware.GenerateToken(user.ID)
    if err != nil {
        return "", nil, err
    }

    return token, user, nil
}