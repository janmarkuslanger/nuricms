package service

import (
	"mime/multipart"
	"os"
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

func (s *AssetService) Create(asset *model.Asset) error {
	return s.repo.Create(asset)
}

func (s *AssetService) Save(asset *model.Asset) error {
	return s.repo.Save(asset)
}

func (s *AssetService) FindByID(id uint) (*model.Asset, error) {
	return s.repo.GetOneByID(id)
}

func (s *AssetService) DeleteByID(id uint) error {
	asset, err := s.repo.GetOneByID(id)

	if err != nil {
		return err
	}

	err = s.repo.Delete(asset)

	if err != nil {
		return err
	}

	err = os.Remove(asset.Path)
	return err
}

func (s *AssetService) GetAll() ([]model.Asset, error) {
	return s.repo.GetAll()
}
