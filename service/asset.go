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
	repos *repository.Set
}

func NewAssetService(repos *repository.Set) *AssetService {
	return &AssetService{
		repos: repos,
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
	return s.repos.Asset.Create(asset)
}

func (s *AssetService) Save(asset *model.Asset) error {
	return s.repos.Asset.Save(asset)
}

func (s *AssetService) FindByID(id uint) (*model.Asset, error) {
	return s.repos.Asset.FindByID(id)
}

func (s *AssetService) DeleteByID(id uint) error {
	asset, err := s.repos.Asset.FindByID(id)

	if err != nil {
		return err
	}

	err = s.repos.Asset.Delete(asset)

	if err != nil {
		return err
	}

	err = os.Remove(asset.Path)
	return err
}

func (s *AssetService) GetAll() ([]model.Asset, error) {
	return s.repos.Asset.List()
}
