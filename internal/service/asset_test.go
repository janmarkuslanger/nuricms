package service_test

import (
	"bytes"
	"errors"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/server"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/janmarkuslanger/nuricms/testutils/mockservices"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAssetService_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewAssetService(repos, &mockservices.MockFileOps{})
	a := &model.Asset{Name: "A", Path: "p"}
	err := svc.Create(a)
	assert.NoError(t, err)
	assert.NotZero(t, a.ID)
}

func TestAssetService_Save(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewAssetService(repos, &mockservices.MockFileOps{})
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
	svc := service.NewAssetService(repos, &mockservices.MockFileOps{})
	a := &model.Asset{Name: "C", Path: "p3"}
	svc.Create(a)
	got, err := svc.FindByID(a.ID)
	assert.NoError(t, err)
	assert.Equal(t, a.ID, got.ID)
}

func TestAssetService_List(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewAssetService(repos, &mockservices.MockFileOps{})
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
	svc := service.NewAssetService(repos, &mockservices.MockFileOps{})
	a := &model.Asset{Name: "D", Path: "to_delete.txt"}
	svc.Create(a)
	os.WriteFile(a.Path, []byte("x"), 0644)
	err := svc.DeleteByID(a.ID)
	assert.NoError(t, err)
	_, err = svc.FindByID(a.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.NoFileExists(t, a.Path)
}

type readSeekCloser struct {
	*bytes.Reader
}

func (r *readSeekCloser) Close() error { return nil }

type brokenReader struct{}

func (b *brokenReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}

func (b *brokenReader) ReadAt(p []byte, off int64) (int, error) {
	return 0, errors.New("read error")
}

func (b *brokenReader) Close() error { return nil }

func (b *brokenReader) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

type brokenFileHeader struct{}

func (b *brokenFileHeader) Open() (multipart.File, error) {
	return nil, errors.New("open fail")
}

func (b *brokenFileHeader) Filename() string {
	return "fail.txt"
}

type copyFailHeader struct{}

func (c *copyFailHeader) Open() (multipart.File, error) {
	return &brokenReader{}, nil
}

func (c *copyFailHeader) Filename() string {
	return "broken.txt"
}

func createMultipartFileHeader(t *testing.T, filename string, content []byte) (*multipart.FileHeader, string) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filename)
	assert.NoError(t, err)
	_, err = part.Write(content)
	assert.NoError(t, err)
	writer.Close()

	req := bytes.NewReader(buf.Bytes())
	mr := multipart.NewReader(req, writer.Boundary())
	form, err := mr.ReadForm(1024 * 1024)
	assert.NoError(t, err)

	files := form.File["file"]
	assert.NotEmpty(t, files)
	return files[0], files[0].Filename
}

func Test_UploadFile_Success(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "uploaded")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	mockFS := &mockservices.MockFileOps{Created: tmpFile}
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewAssetService(repos, mockFS)

	header, filename := createMultipartFileHeader(t, "test.txt", []byte("hello"))
	path, err := svc.UploadFile(server.Context{}, header, filename)

	assert.NoError(t, err)
	assert.Contains(t, path, filepath.Join("public", "assets", "test.txt"))
}

func Test_UploadFile_OpenFails(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewAssetService(repos, &mockservices.MockFileOps{})

	header := &brokenFileHeader{}
	_, err := svc.UploadFile(server.Context{}, header, header.Filename())

	assert.EqualError(t, err, "open fail")
}

func Test_UploadFile_MkdirFails(t *testing.T) {
	header, filename := createMultipartFileHeader(t, "fail.txt", []byte("x"))
	mockFS := &mockservices.MockFileOps{MkdirErr: errors.New("mkdir fail")}
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewAssetService(repos, mockFS)

	_, err := svc.UploadFile(server.Context{}, header, filename)
	assert.EqualError(t, err, "mkdir fail")
}

func Test_UploadFile_CreateFails(t *testing.T) {
	header, filename := createMultipartFileHeader(t, "fail.txt", []byte("x"))
	mockFS := &mockservices.MockFileOps{CreateErr: errors.New("create fail")}
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewAssetService(repos, mockFS)

	_, err := svc.UploadFile(server.Context{}, header, filename)
	assert.EqualError(t, err, "create fail")
}

func Test_UploadFile_CopyFails(t *testing.T) {
	tmpFile, _ := os.CreateTemp("", "copyfail")
	defer os.Remove(tmpFile.Name())

	mockFS := &mockservices.MockFileOps{Created: tmpFile}
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewAssetService(repos, mockFS)

	header := &copyFailHeader{}
	_, err := svc.UploadFile(server.Context{}, header, header.Filename())

	assert.EqualError(t, err, "read error")
}
