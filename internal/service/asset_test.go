package service_test

import (
	"bytes"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUploadFile_Success(t *testing.T) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "testfile.txt")
	assert.NoError(t, err)
	part.Write([]byte("test file content"))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	ctx := server.Context{Writer: w, Request: req}

	err = req.ParseMultipartForm(10 << 20)
	assert.NoError(t, err)

	_, header, err := req.FormFile("file")
	assert.NoError(t, err)

	svc := service.NewAssetService(nil)
	path, err := svc.UploadFile(ctx, header)
	assert.NoError(t, err)
	assert.FileExists(t, path)

	defer os.Remove(path)
}

func TestAssetService_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewAssetService(repos)
	a := &model.Asset{Name: "A", Path: "p"}
	err := svc.Create(a)
	assert.NoError(t, err)
	assert.NotZero(t, a.ID)
}

func TestAssetService_Save(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewAssetService(repos)
	a := &model.Asset{Name: "B", Path: "p2"}
	svc.Create(a)
	a.Name = "B2"
	err := svc.Save(a)
	assert.NoError(t, err)
	got, _ := svc.FindByID(a.ID)
	assert.Equal(t, "B2", got.Name)
}

func TestAssetService_FindByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewAssetService(repos)
	a := &model.Asset{Name: "C", Path: "p3"}
	svc.Create(a)
	got, err := svc.FindByID(a.ID)
	assert.NoError(t, err)
	assert.Equal(t, a.ID, got.ID)
}

func TestAssetService_List(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewAssetService(repos)
	for i := 0; i < 3; i++ {
		svc.Create(&model.Asset{Name: "L", Path: "p"})
	}
	list, total, err := svc.List(1, 2)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, list, 2)
}

func TestAssetService_DeleteByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewAssetService(repos)
	a := &model.Asset{Name: "D", Path: "to_delete.txt"}
	svc.Create(a)
	os.WriteFile(a.Path, []byte("x"), 0644)
	err := svc.DeleteByID(a.ID)
	assert.NoError(t, err)
	_, err = svc.FindByID(a.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoFileExists(t, a.Path)
}
