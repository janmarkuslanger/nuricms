package service

import "github.com/janmarkuslanger/nuricms/internal/repository"

type FieldOptionService interface {
}

type fieldOptionService struct {
	repos *repository.Set
}

func NewFieldOptionService(repos *repository.Set) FieldOptionService {
	return &fieldOptionService{repos: repos}
}
