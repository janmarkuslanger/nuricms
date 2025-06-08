package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/model"
)

type ApikeyService struct {
	repos *repository.Set
}

func NewApikeyService(repos *repository.Set) *ApikeyService {
	return &ApikeyService{
		repos: repos,
	}
}

func (s *ApikeyService) Create(name string, ttl time.Duration) (token string, err error) {
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

func (s *ApikeyService) Validate(token string) error {
	apikey, err := s.repos.Apikey.FindByToken(token)
	if err != nil {
		return errors.New("invalid api key")
	}
	if apikey.ExpiresAt != nil && time.Now().After(*apikey.ExpiresAt) {
		return errors.New("api key expired")
	}
	return nil
}

func (s *ApikeyService) List(page, pageSize int) ([]model.Apikey, int64, error) {
	return s.repos.Apikey.List(page, pageSize)
}

func (s *ApikeyService) FindByID(id uint) (*model.Apikey, error) {
	return s.repos.Apikey.FindByID(id)
}

func (s *ApikeyService) DeleteByID(id uint) error {
	apikey, err := s.repos.Apikey.FindByID(id)
	if err != nil {
		return err
	}

	return s.repos.Apikey.Delete(apikey)
}
