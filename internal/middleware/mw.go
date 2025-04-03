package middleware

import (
	"net/http"
	"strings"

	"fileupbackendv2/internal/auth"
	"fileupbackendv2/internal/models"
	"fileupbackendv2/internal/utils"
	"fileupbackendv2/pkg/response"
)

var ExemptPaths = map[string]struct{}{"/files/download/": {}, "/auth/login/": {}, "/auth/signup/": {}, "/auth/generate-api-key/": {}}

func CheckExemptPath(path string) bool {
	for exemptPath := range ExemptPaths {
		if len(path) > len(exemptPath) {
			if path[:len(exemptPath)] == exemptPath {
				return true
			}
		} else if path == exemptPath {
			return true
		}
	}
	return false
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exempt := CheckExemptPath(r.URL.Path)
		if exempt {
			next.ServeHTTP(w, r)
			return
		}
		apiKeyHeader := r.Header.Get("X-API-Key")
		tokenHeader := r.Header.Get("Authorization")
		tokenHeaderSplit := strings.Split(tokenHeader, " ")
		if len(tokenHeaderSplit) == 2 {
			tokenHeader = tokenHeaderSplit[1]
		}
		if apiKeyHeader == "" && tokenHeader == "" {
			response.UnauthorizedResponse(w, "Missing API Key or Token")
			return
		}
		var user models.User
		var err error
		if tokenHeader == "" {
			user, err = auth.GetUserByKey(apiKeyHeader)
			if err != nil {
				response.UnauthorizedResponse(w, "Invalid API Key")
				return
			}
		} else {
			claims, err := utils.VerifyToken(tokenHeader)
			if err != nil {
				response.UnauthorizedResponse(w, "Invalid Token")
				return
			}
			user, err = models.GetUserById(claims.UserID)
			if err != nil {
				response.InternalServerErrorResponse(w, err.Error())
				return
			}
		}
		r.Header.Set("X-API-Key", apiKeyHeader)
		r.Header.Set("username", user.Username)
		r.Header.Set("userId", user.ID.String())
		next.ServeHTTP(w, r)
	})
}
