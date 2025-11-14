package middleware

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/okunix/prservice/internal/app/config"
	"github.com/okunix/prservice/internal/pkg/models"
)

const BearerPrefix = "Bearer "

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		adminToken := config.GetConfig().AdminToken

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, BearerPrefix) {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrNotFound)
			return
		}

		token := strings.Replace(authHeader, BearerPrefix, "", 1)

		if subtle.ConstantTimeCompare([]byte(token), []byte(adminToken)) == 0 {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrNotFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
