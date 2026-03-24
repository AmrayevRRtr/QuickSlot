package service

import (
	"QuickSlot/internal/model"
	"QuickSlot/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.UserRepository
}

func NewAuthService(r repository.UserRepository) *AuthService {
	return &AuthService{repo: r}
}

func (s *AuthService) Register(email, password string) (int64, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := &model.User{
		Email:    email,
		Password: string(hash),
		Role:     "USER",
	}

	return s.repo.CreateUser(user)
}

func (s *AuthService) Login(email, password string) (*model.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, err

}
