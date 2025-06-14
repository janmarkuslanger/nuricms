package service

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/testutils"
)

func TestUserService_Create_ValidRole(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := NewUserService(repos, []byte("secret"))

	user, err := svc.Create("u@example.com", "pass123", model.RoleAdmin)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "u@example.com", user.Email)
	assert.Equal(t, model.RoleAdmin, user.Role)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("pass123")))
}

func TestUserService_Create_InvalidRole(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := NewUserService(repos, []byte("secret"))

	user, err := svc.Create("u2@example.com", "pass", model.Role("unknown"))
	assert.Nil(t, user)
	assert.EqualError(t, err, "Not a valid role")
}

func TestUserService_ListAndDeleteByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := NewUserService(repos, []byte("secret"))

	u1, _ := svc.Create("a@e.com", "p", model.RoleEditor)
	svc.Create("ab@e.com", "p", model.RoleEditor)

	list, total, err := svc.List(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, list, 2)

	err = svc.DeleteByID(u1.ID)
	assert.NoError(t, err)
	_, err = repos.User.FindByID(u1.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestUserService_FindSaveDelete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := NewUserService(repos, []byte("secret"))

	u, _ := svc.Create("c@e.com", "pw", model.RoleAdmin)
	found, err := svc.FindByID(u.ID)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, found.ID)

	found.Role = model.RoleEditor
	err = svc.Save(found)
	assert.NoError(t, err)
	reloaded, _ := svc.FindByID(u.ID)
	assert.Equal(t, model.RoleEditor, reloaded.Role)

	err = svc.Delete(reloaded)
	assert.NoError(t, err)
	_, err = svc.FindByID(u.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestUserService_LoginUser_Success(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	secret := []byte("mysecret")
	svc := NewUserService(repos, secret)

	u, _ := svc.Create("x@e.com", "mypw", model.RoleEditor)

	tokenStr, err := svc.LoginUser("x@e.com", "mypw")
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenStr)

	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		return secret, nil
	})
	assert.NoError(t, err)
	assert.True(t, tok.Valid)
	claims := tok.Claims.(jwt.MapClaims)
	assert.EqualValues(t, u.ID, claims["sub"])
	assert.Equal(t, "x@e.com", claims["email"])
	assert.Equal(t, string(model.RoleEditor), claims["role"])
}

func TestUserService_LoginUser_Failure(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := NewUserService(repos, []byte("secret"))

	_, err := svc.LoginUser("no@e.com", "pw")
	assert.Error(t, err)

	svc.Create("y@e.com", "pw2", model.RoleAdmin)
	_, err = svc.LoginUser("y@e.com", "wrong")
	assert.Error(t, err)
}

func TestUserService_ValidateJWT(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	secret := []byte("abc123")
	svc := NewUserService(repos, secret)

	u, _ := svc.Create("z@e.com", "pw3", model.RoleEditor)
	tokenStr, _ := svc.LoginUser("z@e.com", "pw3")

	uid, email, role, err := svc.ValidateJWT(tokenStr)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, uid)
	assert.Equal(t, "z@e.com", email)
	assert.Equal(t, model.RoleEditor, role)
}

func TestUserService_ValidateJWT_Invalid(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := NewUserService(repos, []byte("secret"))

	_, _, _, err := svc.ValidateJWT("notatoken")
	assert.Error(t, err)
}
