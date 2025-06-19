package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
)

type ApikeyService interface {
	List(page, pageSize int) ([]model.Apikey, int64, error)
	CreateToken(name string, ttl time.Duration) (string, error)
	FindByID(id uint) (*model.Apikey, error)
	DeleteByID(id uint) error
	Validate(token string) error
}

type apikeyService struct {
	repos *repository.Set
}

func NewApikeyService(repos *repository.Set) *apikeyService {
	return &apikeyService{
		repos: repos,
	}
}

func (s *apikeyService) CreateToken(name string, ttl time.Duration) (token string, err error) {
	if name == "" {
		return "", errors.New("no name given")
	}

	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", err
	}
	token = hex.EncodeToString(b)

	apiKey := &model.Apikey{
		Name:  name,
		Token: token,
	}
	if ttl > 0 {
		exp := time.Now().Add(ttl)
		apiKey.ExpiresAt = &exp
	}

	if err := s.repos.Apikey.Create(apiKey); err != nil {
		return "", err
	}
	return token, nil
}

func (s *apikeyService) Validate(token string) error {
	apikey, err := s.repos.Apikey.FindByToken(token)
	if err != nil {
		return errors.New("invalid api key")
	}
	if apikey.ExpiresAt != nil && time.Now().After(*apikey.ExpiresAt) {
		return errors.New("api key expired")
	}
	return nil
}

func (s *apikeyService) List(page, pageSize int) ([]model.Apikey, int64, error) {
	return s.repos.Apikey.List(page, pageSize)
}

func (s *apikeyService) FindByID(id uint) (*model.Apikey, error) {
	return s.repos.Apikey.FindByID(id)
}

func (s *apikeyService) DeleteByID(id uint) error {
	apikey, err := s.repos.Apikey.FindByID(id)
	if err != nil {
		return err
	}

	return s.repos.Apikey.Delete(apikey)
}
