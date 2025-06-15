package service

import (
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
)

type AssetService interface {
	List(page, pageSize int) ([]model.Asset, int64, error)
	DeleteByID(id uint) error
	Save(asset *model.Asset) error
	Create(asset *model.Asset) error
	UploadFile(c *gin.Context, file *multipart.FileHeader) (string, error)
	FindByID(id uint) (*model.Asset, error)
}

type assetService struct {
	repos *repository.Set
}

func NewAssetService(repos *repository.Set) AssetService {
	return &assetService{
		repos: repos,
	}
}

func (s *assetService) UploadFile(c *gin.Context, file *multipart.FileHeader) (string, error) {
	savePath := filepath.Join("public", "assets", file.Filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return savePath, err
	}

	return savePath, nil
}

func (s *assetService) Create(asset *model.Asset) error {
	return s.repos.Asset.Create(asset)
}

func (s *assetService) Save(asset *model.Asset) error {
	return s.repos.Asset.Save(asset)
}

func (s *assetService) FindByID(id uint) (*model.Asset, error) {
	return s.repos.Asset.FindByID(id)
}

func (s *assetService) DeleteByID(id uint) error {
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

func (s *assetService) List(page, pageSize int) ([]model.Asset, int64, error) {
	return s.repos.Asset.List(page, pageSize)
}
