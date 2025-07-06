package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/janmarkuslanger/nuricms/internal/model"
)

type jwtValidator interface {
	ValidateJWT(token string) (uint, string, model.Role, error)
}

func Userauth(jwtService jwtValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenCookie, err := r.Cookie("auth_token")
			var token string
			if err == nil {
				token = tokenCookie.Value
			} else {
				authHeader := r.Header.Get("Authorization")
				if !strings.HasPrefix(authHeader, "Bearer ") {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
					return
				}
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}

			uid, email, role, err := jwtService.ValidateJWT(token)
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, uid)
			ctx = context.WithValue(ctx, UserEmailKey, email)
			ctx = context.WithValue(ctx, UserRoleKey, role)

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func Roleauth(allowed ...model.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleVal := r.Context().Value(UserRoleKey)
			role, ok := roleVal.(model.Role)
			if !ok {
				http.Error(w, "missing role", http.StatusForbidden)
				return
			}

			for _, want := range allowed {
				if role == want {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "insufficient permissions", http.StatusForbidden)
		})
	}
}
