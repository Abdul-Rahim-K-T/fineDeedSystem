package middleware

import (
	"context"
	"fineDeedSystem/admin-service/pkg/jwt"
	"log"

	"net/http"
	"strings"
)

// Define a custom type for context keys
type contextKey string

const userContextKey contextKey = "user"

// AuthMiddleware validates the JWT token
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string

		// Check Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			} else {
				log.Println("Invalid token format")
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}
		} else {
			// Check Cookie
			cookie, err := r.Cookie("token")
			if err != nil {
				log.Println("Missing token")
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}
			token = cookie.Value
		}

		log.Printf("Extracted Token: %s", token) // Debug log for extracted token

		claims, err := jwt.ValidateToken(token)
		if err != nil {
			log.Printf("Token validation error: %v", err) // Log the error for debugging
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Set the claims in the request context if needed
		ctx := context.WithValue(r.Context(), userContextKey, claims)
		r = r.WithContext(ctx)

		next(w, r)
	}
}
