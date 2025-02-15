package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwt_key string = os.Getenv("JWT_KEY")

func CreateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"role":     "user",
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iss":      "NexusMeet.server",
		"aud":      "NexusMeet.client",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwt_key))
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

func GetPayloadFromToken(token string) (jwt.MapClaims, error) {
	hmacSecretString := jwt_key
	hmacSecret := []byte(hmacSecretString)

	// Parse the token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid JWT token")
	}
}
