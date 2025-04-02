package middleware

import (
	"net/http"

	"fileupbackendv2/internal/storage"
	"fileupbackendv2/pkg/logging"
	"fileupbackendv2/pkg/response"
)

var ExemptPaths = map[string]struct{}{"/files/download/": {}}

func CheckExemptPath(path string) bool {
	for exemptPath := range ExemptPaths {
		if len(path) > len(exemptPath) {
			if path[:len(exemptPath)] == exemptPath {
				return true
			}
		}
	}
	return false
}

func APIKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exempt := CheckExemptPath(r.URL.Path)
		if exempt {
			next.ServeHTTP(w, r)
			return
		}
		apiKeyHeader := r.Header.Get("X-API-Key")
		logging.AuditLogger.Info().Msgf("API Key: %s", apiKeyHeader)
		if apiKeyHeader == "" {
			response.UnauthorizedResponse(w, "Missing API Key")
			return
		}
		user, err := storage.GetUserByKey(apiKeyHeader)
		if err != nil {
			response.UnauthorizedResponse(w, "Invalid API Key")
			return
		}
		r.Header.Set("X-API-Key", apiKeyHeader)
		r.Header.Set("username", user.Username)
		r.Header.Set("userId", user.ID.String())
		next.ServeHTTP(w, r)
	})
}
