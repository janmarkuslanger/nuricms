package service_test

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/janmarkuslanger/nuricms/internal/dto"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/repository"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
)

func TestCreate_ValidRole(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewUserService(repos, []byte("secret"))

	user, err := svc.Create(dto.UserData{
		Email:    "u@example.com",
		Password: "pass123",
		Role:     string(model.RoleAdmin),
	})
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.Equal(t, "u@example.com", user.Email)
	assert.Equal(t, model.RoleAdmin, user.Role)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("pass123")))
}

func TestCreate_InvalidRole(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewUserService(repos, []byte("secret"))

	user, err := svc.Create(dto.UserData{
		Email:    "u@example.com",
		Password: "pass",
		Role:     "unknown",
	})
	assert.Nil(t, user)
	assert.EqualError(t, err, "not a valid role")
}

func TestListAndDeleteByID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewUserService(repos, []byte("secret"))

	u1, _ := svc.Create(dto.UserData{
		Email:    "a@e.com",
		Password: "p",
		Role:     string(model.RoleEditor),
	})
	svc.Create(dto.UserData{
		Email:    "aaaaaa@e.com",
		Password: "p",
		Role:     string(model.RoleEditor),
	})

	list, total, err := svc.List(1, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, list, 2)

	err = svc.DeleteByID(u1.ID)
	assert.NoError(t, err)
	_, err = repos.User.FindByID(u1.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestFindSaveDelete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewUserService(repos, []byte("secret"))

	u, _ := svc.Create(dto.UserData{
		Email:    "c@e.com",
		Password: "pw",
		Role:     string(model.RoleAdmin),
	})
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

func TestLoginUser_Success(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	secret := []byte("mysecret")
	svc := service.NewUserService(repos, secret)

	u, _ := svc.Create(dto.UserData{
		Email:    "x@e.com",
		Password: "mypw",
		Role:     string(model.RoleEditor),
	})

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

func TestLoginUser_EmptyEmail(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewUserService(repos, []byte("secret"))

	_, err := svc.LoginUser("", "pw")
	assert.Error(t, err)
	assert.EqualError(t, err, "email is empty")
}

func TestLoginUser_EmptyPassword(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewUserService(repos, []byte("secret"))

	_, err := svc.LoginUser("nuri@nuri.com", "")
	assert.Error(t, err)
	assert.EqualError(t, err, "password is empty")
}

func TestLoginUser_Failure(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewUserService(repos, []byte("secret"))

	_, err := svc.LoginUser("no@e.com", "pw")
	assert.Error(t, err)

	svc.Create(dto.UserData{
		Email:    "y@e.com",
		Password: "pw2",
		Role:     string(model.RoleAdmin),
	})
	_, err = svc.LoginUser("y@e.com", "wrong")
	assert.Error(t, err)
}

func TestValidateJWT(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	secret := []byte("abc123")
	svc := service.NewUserService(repos, secret)

	u, _ := svc.Create(dto.UserData{
		Email:    "z@e.com",
		Password: "pw3",
		Role:     string(model.RoleEditor),
	})
	tokenStr, _ := svc.LoginUser("z@e.com", "pw3")

	uid, email, role, err := svc.ValidateJWT(tokenStr)
	assert.NoError(t, err)
	assert.Equal(t, u.ID, uid)
	assert.Equal(t, "z@e.com", email)
	assert.Equal(t, model.RoleEditor, role)
}

func TestValidateJWT_Invalid(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewUserService(repos, []byte("secret"))

	_, _, _, err := svc.ValidateJWT("notatoken")
	assert.Error(t, err)
}

func TestUpdateByID_NoEmail(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewUserService(repos, []byte("secret"))

	svc.Create(dto.UserData{
		Email:    "test",
		Password: "p",
		Role:     string(model.RoleEditor),
	})

	_, err := svc.UpdateByID(1, dto.UserData{
		Password: "p",
		Role:     string(model.RoleEditor),
	})

	assert.Equal(t, err.Error(), "no email given")
}

func TestUpdateByID_NoPw(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewUserService(repos, []byte("secret"))

	svc.Create(dto.UserData{
		Email:    "test",
		Password: "p",
		Role:     string(model.RoleEditor),
	})

	_, err := svc.UpdateByID(1, dto.UserData{
		Email: "EMAIL",
		Role:  string(model.RoleEditor),
	})

	assert.Equal(t, err.Error(), "no password given")
}

func TestUpdateByID_NoRole(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewUserService(repos, []byte("secret"))

	svc.Create(dto.UserData{
		Email:    "test",
		Password: "p",
		Role:     string(model.RoleEditor),
	})

	_, err := svc.UpdateByID(1, dto.UserData{
		Email:    "EMAIL",
		Password: "Test",
	})

	assert.Equal(t, err.Error(), "no role given")
}

func TestUpdateByID_Success(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repos := repository.NewSet(db)
	svc := service.NewUserService(repos, []byte("secret"))

	svc.Create(dto.UserData{
		Email:    "beforeE",
		Password: "beforePw",
		Role:     string(model.RoleEditor),
	})

	u, err := svc.UpdateByID(1, dto.UserData{
		Email:    "afterE",
		Password: "afterPw",
		Role:     string(model.RoleAdmin),
	})

	assert.Equal(t, err, nil)
	assert.Equal(t, u.Email, "afterE")
	assert.Equal(t, u.Password, "afterPw")
	assert.Equal(t, u.Role, model.RoleAdmin)
}
