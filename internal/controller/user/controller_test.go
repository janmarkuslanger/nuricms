package user

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/janmarkuslanger/nuricms/internal/service"
	"github.com/janmarkuslanger/nuricms/testutils"
	"github.com/stretchr/testify/assert"
)

func createMockController(userSvc service.UserService) *Controller {
	gin.SetMode(gin.TestMode)
	svcSet := &service.Set{User: userSvc}
	return NewController(svcSet)
}

func TestUserController_login_Success(t *testing.T) {
	svc := new(testutils.MockUserService)
	svc.On("LoginUser", "test@example.com", "pass123").Return("jwt-token", nil)
	ct := createMockController(svc)

	c, w := testutils.MakePOSTContext("/login", gin.H{
		"email":    "test@example.com",
		"password": "pass123",
	})

	ct.login(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/", w.Header().Get("Location"))
	assert.Contains(t, w.Header().Get("Set-Cookie"), "auth_token=jwt-token")
	svc.AssertExpectations(t)
}

func TestUserController_login_Failed(t *testing.T) {
	svc := new(testutils.MockUserService)
	svc.On("LoginUser", "test@example.com", "wrongpass").Return("", errors.New("invalid credentials"))
	ct := createMockController(svc)

	c, w := testutils.MakePOSTContext("/login", gin.H{
		"email":    "test@example.com",
		"password": "wrongpass",
	})

	ct.login(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/login", w.Header().Get("Location"))
	svc.AssertExpectations(t)
}

func TestUserController_createUser_Success(t *testing.T) {
	svc := new(testutils.MockUserService)
	svc.On("Create", "test@example.com", "secure123", model.RoleEditor).
		Return(&model.User{}, nil)
	ct := createMockController(svc)

	c, w := testutils.MakePOSTContext("/user/create", gin.H{
		"email":    "test@example.com",
		"password": "secure123",
		"role":     string(model.RoleEditor),
	})

	ct.createUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/user", w.Header().Get("Location"))
	svc.AssertExpectations(t)
}

func TestUserController_createUser_Failed(t *testing.T) {
	svc := new(testutils.MockUserService)
	svc.On("Create", "fail@example.com", "123", model.RoleAdmin).
		Return((*model.User)(nil), errors.New("failed"))
	ct := createMockController(svc)

	c, w := testutils.MakePOSTContext("/user/create", gin.H{
		"email":    "fail@example.com",
		"password": "123",
		"role":     string(model.RoleAdmin),
	})

	ct.createUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "/user", w.Header().Get("Location"))
	svc.AssertExpectations(t)
}

func TestUserController_showUser_Success(t *testing.T) {
	svc := new(testutils.MockUserService)
	svc.On("List", 1, 10).Return([]model.User{}, int64(0), nil)
	ct := createMockController(svc)

	teardown := testutils.StubRenderWithLayout()
	defer teardown()

	c, w := testutils.MakeGETContext("/user?page=1&pageSize=10")
	ct.showUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "RENDERED:user/index.tmpl")
	svc.AssertExpectations(t)
}
