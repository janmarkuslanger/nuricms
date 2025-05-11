package service

import (
	"mime/multipart"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/model"
	"github.com/janmarkuslanger/nuricms/repository"
)

type AssetService struct {
	repo *repository.AssetRepository
}

func NewAssetService(repo *repository.AssetRepository) *AssetService {
	return &AssetService{
		repo: repo,
	}
}

func (s *AssetService) UploadFile(c *gin.Context, file *multipart.FileHeader) (string, error) {
	savePath := filepath.Join("public", "assets", file.Filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return savePath, err
	}

	return savePath, nil
}

func (s *AssetService) Create(asset *model.Asset) (*model.Asset, error) {
	_, err := s.repo.Create(asset)
	return asset, err
}
