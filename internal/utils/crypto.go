package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fileupbackendv2/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey []byte

func init() {
	secretKey = []byte(config.SecretKey)
}

func GenerateKey(n int) (string, error) {
	key := make([]byte, n)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	apiKey := base64.RawURLEncoding.EncodeToString(key)
	return apiKey, nil
}
func HashKey(apiKey string) string {
	combined := append([]byte(apiKey), []byte(config.SecretKey)...)
	hash := sha256.New()
	hash.Write(combined)
	hashedBytes := hash.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}

// Custom claims struct to hold user data
type Claims struct {
	Username string `json:"username"`
	UserID   string
	jwt.RegisteredClaims
}

func CreateToken(username string, userId string) (string, error) {
	// Create claims with expiration time (e.g., 24 hours)
	claims := Claims{
		Username: username,
		UserID:   userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*Claims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Verify token is valid and extract claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
