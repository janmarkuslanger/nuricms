package repository

import (
	"errors"
	"fmt"
	"testing"

	"github.com/janmarkuslanger/nuricms/model"
	"gorm.io/gorm"
)

func setupUserRepo(t *testing.T) (*UserRepository, func()) {
	t.Helper()

	db, err := CreateTestDB()
	if err != nil {
		t.Fatalf("CreateTestDB failed: %v", err)
	}

	repo := NewUserRepository(db)

	cleanup := func() {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	return repo, cleanup
}

func TestUserRepository_Create_And_FindByID_And_FindByEmail(t *testing.T) {
	repo, cleanup := setupUserRepo(t)
	defer cleanup()

	u := &model.User{
		Email:    "alice@example.com",
		Password: "password",
	}
	if err := repo.Create(u); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if u.ID == 0 {
		t.Fatalf("Expected ID to be set, got 0")
	}

	foundByID, err := repo.FindByID(u.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if foundByID.Email != u.Email {
		t.Errorf("FindByID: expected email %q, got %q", u.Email, foundByID.Email)
	}

	foundByEmail, err := repo.FindByEmail("alice@example.com")
	if err != nil {
		t.Fatalf("FindByEmail failed: %v", err)
	}
	if foundByEmail.ID != u.ID {
		t.Errorf("FindByEmail: expected ID %d, got %d", u.ID, foundByEmail.ID)
	}

	_, err = repo.FindByEmail("nonexistent@example.com")
	if err == nil {
		t.Errorf("Expected error for non-existing email")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected ErrRecordNotFound, got %v", err)
	}
}

func TestUserRepository_Save_Update_User(t *testing.T) {
	repo, cleanup := setupUserRepo(t)
	defer cleanup()

	u := &model.User{
		Email:    "bob@example.com",
		Password: "pass123",
	}
	if err := repo.Create(u); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	u.Password = "newpassword"
	if err := repo.Save(u); err != nil {
		t.Fatalf("Save (Update) failed: %v", err)
	}

	updated, err := repo.FindByID(u.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if updated.Password != "newpassword" {
		t.Errorf("Save: expected password %q, got %q", "newpassword", updated.Password)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	repo, cleanup := setupUserRepo(t)
	defer cleanup()

	u := &model.User{
		Email:    "chris@example.com",
		Password: "pw",
	}
	if err := repo.Create(u); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	id := u.ID

	if err := repo.Delete(u); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err := repo.FindByID(id)
	if err == nil {
		t.Errorf("Expected error after deleting user")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Expected ErrRecordNotFound, got %v", err)
	}
}

func TestUserRepository_List_Pagination(t *testing.T) {
	repo, cleanup := setupUserRepo(t)
	defer cleanup()

	for i := 1; i <= 5; i++ {
		u := &model.User{
			Email:    fmt.Sprintf("user%02d@example.com", i),
			Password: "pwd",
		}
		if err := repo.Create(u); err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	usersPage1, totalCount, err := repo.List(1, 2)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if totalCount != 5 {
		t.Errorf("List: expected totalCount 5, got %d", totalCount)
	}
	if len(usersPage1) != 2 {
		t.Errorf("Page 1: expected 2 entries, got %d", len(usersPage1))
	}
	if usersPage1[0].Email != "user01@example.com" || usersPage1[1].Email != "user02@example.com" {
		t.Errorf("Page 1: unexpected entries: %+v", usersPage1)
	}

	usersPage3, totalCount3, err := repo.List(3, 2)
	if err != nil {
		t.Fatalf("List (Page 3) failed: %v", err)
	}
	if totalCount3 != 5 {
		t.Errorf("List (Page 3): expected totalCount 5, got %d", totalCount3)
	}
	if len(usersPage3) != 1 {
		t.Errorf("Page 3: expected 1 entry, got %d", len(usersPage3))
	}
	if usersPage3[0].Email != "user05@example.com" {
		t.Errorf("Page 3: expected user05@example.com, got %q", usersPage3[0].Email)
	}
}
