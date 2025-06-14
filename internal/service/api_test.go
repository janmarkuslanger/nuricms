package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/janmarkuslanger/nuricms/internal/db"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestApiService_prepareContent_SimpleAndListCollectionAsset(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	db.DB = testDB
	repos := repository.NewSet(testDB)
	s := NewApiService(repos)

	col1 := &model.Collection{Name: "Col", Alias: "col"}
	repos.Collection.Create(col1)

	fText := &model.Field{Alias: "text", FieldType: model.FieldType("text"), CollectionID: col1.ID}
	repos.Field.Create(fText)
	fList := &model.Field{Alias: "list", FieldType: model.FieldType("text"), CollectionID: col1.ID, IsList: true}
	repos.Field.Create(fList)

	col2 := &model.Collection{Name: "Related", Alias: "related"}
	repos.Collection.Create(col2)
	contentRel := &model.Content{CollectionID: col2.ID}
	repos.Content.Create(contentRel)

	fColl := &model.Field{Alias: "coll", FieldType: model.FieldTypeCollection, CollectionID: col1.ID}
	repos.Field.Create(fColl)
	fAsset := &model.Field{Alias: "asset", FieldType: model.FieldTypeAsset, CollectionID: col1.ID}
	repos.Field.Create(fAsset)

	asset := &model.Asset{Name: "A", Path: "/p"}
	repos.Asset.Create(asset)

	now := time.Now().Truncate(time.Second)
	ce := &model.Content{
		Model:        gorm.Model{ID: 100, CreatedAt: now, UpdatedAt: now},
		CollectionID: col1.ID,
		Collection:   *col1,
	}

	cv1 := model.ContentValue{Model: gorm.Model{ID: 1}, Field: *fText, Value: "hello"}
	cv2 := model.ContentValue{Model: gorm.Model{ID: 2}, Field: *fList, Value: "v1"}
	cv3 := model.ContentValue{Model: gorm.Model{ID: 3}, Field: *fList, Value: "v2"}
	cv4 := model.ContentValue{Model: gorm.Model{ID: 4}, Field: *fColl, Value: fmt.Sprint(contentRel.ID)}
	cv5 := model.ContentValue{Model: gorm.Model{ID: 5}, Field: *fAsset, Value: fmt.Sprint(asset.ID)}

	ce.ContentValues = []model.ContentValue{cv1, cv2, cv3, cv4, cv5}

	resp, err := s.prepareContent(ce)
	assert.NoError(t, err)
	assert.Equal(t, ce.ID, resp.ID)
	assert.Equal(t, now.Unix(), resp.CreatedAt.Unix())
	assert.Equal(t, now.Unix(), resp.UpdatedAt.Unix())
	v := resp.Values["text"].(ContentValueResponse)
	assert.Equal(t, uint(1), v.ID)
	assert.Equal(t, "hello", v.Value)
	listVals := resp.Values["list"].([]any)
	assert.Len(t, listVals, 2)
	collVal := resp.Values["coll"].(ContentValueResponse)
	assert.NotNil(t, collVal.Collection)
	assert.Equal(t, col2.ID, collVal.Collection.ID)
	assetVal := resp.Values["asset"].(ContentValueResponse)
	assert.NotNil(t, assetVal.Asset)
	assert.Equal(t, asset.ID, assetVal.Asset.ID)
	assert.Equal(t, col1.ID, resp.Collection.ID)
	assert.Equal(t, col1.Name, resp.Collection.Name)
	assert.Equal(t, col1.Alias, resp.Collection.Alias)
}

func TestApiService_FindContentByCollectionAlias_FindByID(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	db.DB = testDB
	repos := repository.NewSet(testDB)
	s := NewApiService(repos)

	col := &model.Collection{Name: "ColX", Alias: "colx"}
	repos.Collection.Create(col)
	f := &model.Field{Alias: "f", FieldType: model.FieldType("text"), CollectionID: col.ID}
	repos.Field.Create(f)
	c := &model.Content{CollectionID: col.ID}
	repos.Content.Create(c)
	repos.ContentValue.Create(&model.ContentValue{ContentID: c.ID, FieldID: f.ID, Value: "val"})

	list, err := s.FindContentByCollectionAlias("colx", 0, 10)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	item := list[0]
	assert.Equal(t, c.ID, item.ID)

	single, err := s.FindContentByID(c.ID)
	assert.NoError(t, err)
	assert.Equal(t, c.ID, single.ID)
}

func TestApiService_FindContentByCollectionAndFieldValue(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	db.DB = testDB
	repos := repository.NewSet(testDB)
	s := NewApiService(repos)

	col := &model.Collection{Name: "ColY", Alias: "coly"}
	repos.Collection.Create(col)
	f := &model.Field{Alias: "f2", FieldType: model.FieldType("text"), CollectionID: col.ID}
	repos.Field.Create(f)
	for i := 0; i < 3; i++ {
		c := &model.Content{CollectionID: col.ID}
		repos.Content.Create(c)
		repos.ContentValue.Create(&model.ContentValue{ContentID: c.ID, FieldID: f.ID, Value: "match"})
	}
	items, err := s.FindContentByCollectionAndFieldValue("coly", "f2", "match", 0, 10)
	assert.NoError(t, err)
	assert.Len(t, items, 3)
}

func TestApiService_FindContentByCollectionAlias_Error(t *testing.T) {
	testDB := testutils.SetupTestDB(t)
	db.DB = testDB
	repos := repository.NewSet(testDB)
	s := NewApiService(repos)

	_, err := s.FindContentByCollectionAlias("nonexistent", 0, 10)
	assert.Error(t, err)
}
