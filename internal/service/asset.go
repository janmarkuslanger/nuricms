package service

import (
	"io"
	"path/filepath"

	"github.com/janmarkuslanger/nuricms/internal/fs"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/server"
)

type AssetService interface {
	List(page, pageSize int) ([]model.Asset, int64, error)
	DeleteByID(id uint) error
	Save(asset *model.Asset) error
	Create(asset *model.Asset) error
	UploadFile(ctx server.Context, header fs.FileOpener, filename string) (string, error)
	FindByID(id uint) (*model.Asset, error)
}

type assetService struct {
	repos *repository.Set
	fs    fs.FileOps
}

func NewAssetService(repos *repository.Set, fs fs.FileOps) AssetService {
	return &assetService{
		repos: repos,
		fs:    fs,
	}
}

func (s *assetService) UploadFile(ctx server.Context, header fs.FileOpener, filename string) (string, error) {
	src, err := header.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	savePath := filepath.Join("public", "assets", filename)

	if err := s.fs.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		return "", err
	}

	dst, err := s.fs.Create(savePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
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

	return s.fs.Remove(asset.Path)
}

func (s *assetService) List(page, pageSize int) ([]model.Asset, int64, error) {
	return s.repos.Asset.List(page, pageSize)
}
