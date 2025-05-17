package repository

import (
	"testing"

	"github.com/janmarkuslanger/nuricms/model"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestAssetRepository_CreateAndGetOneByID(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}

	repo := NewAssetRepository(db)
	asset := &model.Asset{}
	err = repo.Create(asset)
	require.NoError(t, err)
	require.NotZero(t, asset.ID)
	fetched, err := repo.FindByID(asset.ID)
	require.NoError(t, err)
	require.Equal(t, asset.ID, fetched.ID)
}

func TestAssetRepository_GetOneByID_NotFound(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewAssetRepository(db)
	_, err = repo.FindByID(999)
	require.Error(t, err)
	require.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestAssetRepository_GetAll(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewAssetRepository(db)
	assets, err := repo.List()
	require.NoError(t, err)
	require.Len(t, assets, 0)
	for i := 0; i < 2; i++ {
		a := &model.Asset{}
		err := repo.Create(a)
		require.NoError(t, err)
	}
	assets, err = repo.List()
	require.NoError(t, err)
	require.Len(t, assets, 2)
}

func TestAssetRepository_Save(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewAssetRepository(db)
	asset := &model.Asset{}
	err = repo.Create(asset)
	require.NoError(t, err)
	err = repo.Save(asset)
	require.NoError(t, err)
	_, err = repo.FindByID(asset.ID)
	require.NoError(t, err)
}

func TestAssetRepository_Delete(t *testing.T) {
	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	repo := NewAssetRepository(db)
	asset := &model.Asset{}
	err = repo.Create(asset)
	require.NoError(t, err)
	err = repo.Delete(asset)
	require.NoError(t, err)
	_, err = repo.FindByID(asset.ID)
	require.Error(t, err)
	require.Equal(t, gorm.ErrRecordNotFound, err)
}
