package middleware

import (
	"encoding/json"
	"net/http"
)

type validator interface {
	Validate(token string) error
}

func ApikeyAuth(keySvc validator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("X-API-Key")
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "missing X-API-Key",
				})
				return
			}

			if err := keySvc.Validate(token); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(map[string]string{
					"error": err.Error(),
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
