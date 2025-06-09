package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/janmarkuslanger/nuricms/internal/model"
	"github.com/stretchr/testify/assert"
)

type fakeUserService struct {
	validTokens map[string]struct {
		uid   uint
		email string
		role  model.Role
	}
}

func (f *fakeUserService) ValidateJWT(token string) (uint, string, model.Role, error) {
	if val, ok := f.validTokens[token]; ok {
		return val.uid, val.email, val.role, nil
	}
	return 0, "", "", errors.New("invalid token")
}

func TestUserauthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fakeSvc := &fakeUserService{
		validTokens: map[string]struct {
			uid   uint
			email string
			role  model.Role
		}{
			"good-token": {uid: 42, email: "test@example.com", role: model.RoleAdmin},
		},
	}

	tests := []struct {
		name             string
		setCookie        bool
		cookieToken      string
		setAuthHeader    bool
		headerValue      string
		expectedStatus   int
		expectedLocation string
		expectUserID     uint
		expectEmail      string
		expectRole       model.Role
		callNext         bool
	}{
		{
			name:             "No cookie, no Authorization header → Redirect",
			setCookie:        false,
			setAuthHeader:    false,
			expectedStatus:   http.StatusSeeOther,
			expectedLocation: "/login",
			callNext:         false,
		},
		{
			name:             "Invalid Authorization header (no Bearer prefix) → Redirect",
			setCookie:        false,
			setAuthHeader:    true,
			headerValue:      "SomethingElse xyz",
			expectedStatus:   http.StatusSeeOther,
			expectedLocation: "/login",
			callNext:         false,
		},
		{
			name:             "Bearer header present, invalid token → Redirect",
			setCookie:        false,
			setAuthHeader:    true,
			headerValue:      "Bearer bad-token",
			expectedStatus:   http.StatusSeeOther,
			expectedLocation: "/login",
			callNext:         false,
		},
		{
			name:           "Bearer header present, valid token → Next called",
			setCookie:      false,
			setAuthHeader:  true,
			headerValue:    "Bearer good-token",
			expectedStatus: http.StatusOK,
			expectUserID:   42,
			expectEmail:    "test@example.com",
			expectRole:     model.RoleAdmin,
			callNext:       true,
		},
		{
			name:           "Cookie present, valid token → Next called",
			setCookie:      true,
			cookieToken:    "good-token",
			setAuthHeader:  false,
			expectedStatus: http.StatusOK,
			expectUserID:   42,
			expectEmail:    "test@example.com",
			expectRole:     model.RoleAdmin,
			callNext:       true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.Use(Userauth(fakeSvc))
			router.GET("/test", func(c *gin.Context) {
				id, _ := c.Get("userID")
				email, _ := c.Get("userEmail")
				role, _ := c.Get("userRole")
				c.JSON(http.StatusOK, gin.H{
					"userID":    id,
					"userEmail": email,
					"userRole":  role,
				})
			})

			req := httptest.NewRequest("GET", "/test", nil)
			if tc.setCookie {
				req.AddCookie(&http.Cookie{Name: "auth_token", Value: tc.cookieToken})
			}
			if tc.setAuthHeader {
				req.Header.Set("Authorization", tc.headerValue)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if tc.callNext {
				assert.Equal(t, http.StatusOK, w.Code)
				var body map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &body)
				assert.NoError(t, err)
				assert.Equal(t, float64(tc.expectUserID), body["userID"])
				assert.Equal(t, tc.expectEmail, body["userEmail"])
				assert.Equal(t, string(tc.expectRole), body["userRole"])
			} else {
				assert.Equal(t, tc.expectedStatus, w.Code)
				loc := w.Header().Get("Location")
				assert.Equal(t, tc.expectedLocation, loc)
			}
		})
	}
}

func TestRoleauthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userRole       model.Role
		allowedRoles   []model.Role
		expectedStatus int
		responseBody   string
		callNext       bool
	}{
		{
			name:           "RoleAdmin allowed → Next called",
			userRole:       model.RoleAdmin,
			allowedRoles:   []model.Role{model.RoleAdmin, model.RoleEditor},
			expectedStatus: http.StatusOK,
			callNext:       true,
		},
		{
			name:           "RoleEditor not allowed (only Admin) → 403 Forbidden",
			userRole:       model.RoleEditor,
			allowedRoles:   []model.Role{model.RoleAdmin},
			expectedStatus: http.StatusForbidden,
			responseBody:   `{"error":"insufficient permissions"}`,
			callNext:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router := gin.New()
			router.Use(func(c *gin.Context) {
				c.Set("userRole", tc.userRole)
				c.Next()
			})
			router.GET("/test", Roleauth(tc.allowedRoles...), func(c *gin.Context) {
				c.String(http.StatusOK, "ROLE_OK")
			})

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			if tc.callNext {
				assert.Equal(t, "ROLE_OK", w.Body.String())
			} else {
				assert.JSONEq(t, tc.responseBody, w.Body.String())
			}
		})
	}
}
