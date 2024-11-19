package middleware

import (
	"fineDeedSystem/applicant-service/internal/usecase"
	"net/http"
	"strings"
)

func ApplicantAuthMiddleware(uc *usecase.ApplicantUsecase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token string

			// First check the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				// Remove "Bearer " if present in the Authorization header
				if strings.HasPrefix(authHeader, "Bearer ") {
					token = authHeader[len("Bearer "):]
				}
			}

			// If token not found in header, check for token in cookies
			if token == "" {
				cookie, err := r.Cookie("auth_token")
				if err == nil {
					token = cookie.Value
				}
			}

			if token == "" {
				http.Error(w, "No token provided", http.StatusUnauthorized)
				return
			}

			// Check if the token is valid
			if !uc.IsTokenValid(token) {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Proceed with the next handler
			next.ServeHTTP(w, r)
		})
	}
}
