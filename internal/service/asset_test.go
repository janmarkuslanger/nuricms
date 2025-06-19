package service

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAssetService_UploadFile(t *testing.T) {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	fw, _ := writer.CreateFormFile("file", "upload.txt")
	fw.Write([]byte("content"))
	writer.Close()
	ctx.Request = &http.Request{
		Method: "POST",
		Header: http.Header{"Content-Type": {writer.FormDataContentType()}},
		Body:   io.NopCloser(&buf),
	}
	file, err := ctx.FormFile("file")
	assert.NoError(t, err)
	svc := NewAssetService(&repository.Set{})
	os.MkdirAll("public/assets", 0755)
	defer os.RemoveAll("public")
	path, err := svc.UploadFile(ctx, file)
	assert.NoError(t, err)
	assert.FileExists(t, path)
}

func TestAssetService_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := NewAssetService(repos)
	a := &model.Asset{Name: "A", Path: "p"}
	err := svc.Create(a)
	assert.NoError(t, err)
	assert.NotZero(t, a.ID)
}

func TestAssetService_Save(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := NewAssetService(repos)
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
	svc := NewAssetService(repos)
	a := &model.Asset{Name: "C", Path: "p3"}
	svc.Create(a)
	got, err := svc.FindByID(a.ID)
	assert.NoError(t, err)
	assert.Equal(t, a.ID, got.ID)
}

func TestAssetService_List(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := NewAssetService(repos)
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
	svc := NewAssetService(repos)
	a := &model.Asset{Name: "D", Path: "to_delete.txt"}
	svc.Create(a)
	os.WriteFile(a.Path, []byte("x"), 0644)
	err := svc.DeleteByID(a.ID)
	assert.NoError(t, err)
	_, err = svc.FindByID(a.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoFileExists(t, a.Path)
}
