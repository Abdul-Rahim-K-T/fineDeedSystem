package middleware

import (
	"context"
	"fineDeedSystem/employer-service/internal/usecase"
	"fineDeedSystem/employer-service/pkg/constants"
	"fineDeedSystem/employer-service/pkg/jwt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func EmployerAuthMiddleware(usecase *usecase.EmployerUsecase) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var tokenString string

			// First, try to get the token from the Authorization header
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				cookie, err := r.Cookie("jwt_token")
				if err != nil || cookie.Value == "" {
					http.Error(w, "No token provided", http.StatusUnauthorized)
					return
				}
				tokenString = cookie.Value
			}

			// Log the extracted token for debugging purposes
			log.Printf("Authorization Haeder: %s", authHeader)
			log.Printf("Extracted Token from Cookie: %s", tokenString)
			log.Printf("Received Token: %s", tokenString)

			if tokenString == "" {
				http.Error(w, "No token provided", http.StatusUnauthorized) // Debugggggggggggg
				return
			}

			claims, err := jwt.ParseToken(tokenString)
			if err != nil {
				log.Printf("Token parsing Error: %v", err) // Debugging statement
				http.Error(w, "Invalid token or not an employer", http.StatusUnauthorized)
				return
			}

			log.Printf("Token Claims: %+v", claims) // Debugging statement

			if !claims.IsEmployer {
				http.Error(w, "Invalid token or not an employer", http.StatusUnauthorized)
				return
			}

			// Check if the token is blacklisted
			isBlacklisted, err := usecase.IsTokenBlacklisted(tokenString)
			if err != nil {
				log.Println("Error checking token blacklist:", err)
				http.Error(w, "Error checking token blacklist", http.StatusInternalServerError)
				return
			}
			if isBlacklisted {
				http.Error(w, "Token is blacklisted", http.StatusUnauthorized)
				return
			}

			// Pass employer ID to the next handler
			ctx := context.WithValue(r.Context(), constants.UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
