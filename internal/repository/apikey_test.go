package repository

import (
	"errors"
	"fmt"
	"testing"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"gorm.io/gorm"
)

func setupApikeyRepo(t *testing.T) (*ApikeyRepository, func()) {
	t.Helper()

	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("CreateTestDB failed: %v", err)
	}

	repo := NewApikeyRepository(db)

	cleanup := func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	return repo, cleanup
}

func TestApikeyRepository_Create_And_FindByID_And_FindByToken(t *testing.T) {
	repo, cleanup := setupApikeyRepo(t)
	defer cleanup()

	key := &model.Apikey{
		Token: "token123",
	}
	if err := repo.Create(key); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if key.ID == 0 {
		t.Fatalf("Expected ID to be set, got 0")
	}

	foundByID, err := repo.FindByID(key.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if foundByID.Token != key.Token {
		t.Errorf("FindByID: expected token %q, got %q", key.Token, foundByID.Token)
	}

	foundByToken, err := repo.FindByToken("token123")
	if err != nil {
		t.Fatalf("FindByToken failed: %v", err)
	}
	if foundByToken.ID != key.ID {
		t.Errorf("FindByToken: expected ID %d, got %d", key.ID, foundByToken.ID)
	}

	_, err = repo.FindByToken("nonexistent")
	if err == nil {
		t.Errorf("Expected error for non-existing token")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected ErrRecordNotFound, got %v", err)
	}
}

func TestApikeyRepository_Delete(t *testing.T) {
	repo, cleanup := setupApikeyRepo(t)
	defer cleanup()

	key := &model.Apikey{
		Token: "todelete",
	}
	if err := repo.Create(key); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	id := key.ID

	if err := repo.Delete(key); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := repo.FindByID(id)
	if err == nil {
		t.Errorf("Expected error after deleting apikey")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected ErrRecordNotFound, got %v", err)
	}
}

func TestApikeyRepository_List_Pagination(t *testing.T) {
	repo, cleanup := setupApikeyRepo(t)
	defer cleanup()

	for i := 1; i <= 5; i++ {
		key := &model.Apikey{
			Token: fmt.Sprintf("token%02d", i),
		}
		if err := repo.Create(key); err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	keysPage1, totalCount, err := repo.List(1, 2)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if totalCount != 5 {
		t.Errorf("List: expected totalCount 5, got %d", totalCount)
	}
	if len(keysPage1) != 2 {
		t.Errorf("Page 1: expected 2 entries, got %d", len(keysPage1))
	}
	if keysPage1[0].Token != "token01" || keysPage1[1].Token != "token02" {
		t.Errorf("Page 1: unexpected entries: %+v", keysPage1)
	}

	keysPage3, totalCount3, err := repo.List(3, 2)
	if err != nil {
		t.Fatalf("List (Page 3) failed: %v", err)
	}
	if totalCount3 != 5 {
		t.Errorf("List (Page 3): expected totalCount 5, got %d", totalCount3)
	}
	if len(keysPage3) != 1 {
		t.Errorf("Page 3: expected 1 entry, got %d", len(keysPage3))
	}
	if keysPage3[0].Token != "token05" {
		t.Errorf("Page 3: expected token05, got %q", keysPage3[0].Token)
	}
}
