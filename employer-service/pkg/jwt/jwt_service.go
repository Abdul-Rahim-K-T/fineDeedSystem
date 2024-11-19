package jwt

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims struct to encode JWT token
type Claims struct {
	Username   string `json:"username"`
	IsEmployer bool   `json:"is_employer"`
	jwt.StandardClaims
}
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// GenerateToken generates a JWT token for a given username and role
func GenerateToken(username string, isEmployer bool) (string, error) {
	// Define the expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour).Unix()

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username:   username,
		IsEmployer: isEmployer,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken parses a JWT token and returns the claims if the token is valid
func ParseToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	log.Println("Tracking at jwt_service.go")

	if tokenStr == "" {
		return nil, errors.New("empty token string")
	}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		log.Printf("Error parsing token(jwt_service.go): %v", err)
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		log.Println("Token is invalid")
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
