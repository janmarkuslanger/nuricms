package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
)

type ApikeyService interface {
	List(page, pageSize int) ([]model.Apikey, int64, error)
	Create(dto dto.ApikeyData) (*model.Apikey, error)
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

func (s apikeyService) Create(dto dto.ApikeyData) (apiKey *model.Apikey, err error) {
	if dto.Name == "" {
		return apiKey, errors.New("no name given")
	}

	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return apiKey, err
	}
	token := hex.EncodeToString(b)

	apiKey = &model.Apikey{
		Name:  dto.Name,
		Token: token,
	}

	if err := s.repos.Apikey.Create(apiKey); err != nil {
		return nil, err
	}
	return apiKey, nil
}

func (s apikeyService) Validate(token string) error {
	apikey, err := s.repos.Apikey.FindByToken(token)
	if err != nil {
		return errors.New("invalid api key")
	}
	if apikey.ExpiresAt != nil && time.Now().After(*apikey.ExpiresAt) {
		return errors.New("api key expired")
	}
	return nil
}

func (s apikeyService) List(page, pageSize int) ([]model.Apikey, int64, error) {
	return s.repos.Apikey.List(page, pageSize)
}

func (s apikeyService) FindByID(id uint) (*model.Apikey, error) {
	return s.repos.Apikey.FindByID(id)
}

func (s apikeyService) DeleteByID(id uint) error {
	apikey, err := s.repos.Apikey.FindByID(id)
	if err != nil {
		return err
	}

	return s.repos.Apikey.Delete(apikey)
}
