// PATH: go-auth/utils/GenerateHashPassword.go

package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"kplc-outage-app/models"
	"os"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func CompareHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// func GenerateCryptoHash(data string) (string, error) {
// 	hashSHA256 := sha256.New()
// 	hashSHA256.Write([]byte(data))
// 	hashSHA256Sum := hashSHA256.Sum(nil)
// 	bytes, err := hex.EncodeToString(hashSHA256Sum)
// 	return string(bytes), err
// }

func GenerateCryptoHash(data string) string {
	// Create a new SHA-256 hash instance
	hashSHA256 := sha256.New()

	// Write data to the hash
	hashSHA256.Write([]byte(data))

	// Get the final hash as a byte slice
	hashSHA256Sum := hashSHA256.Sum(nil)

	// Encode the hash byte slice to a hex string and return it
	return hex.EncodeToString(hashSHA256Sum)
}

// func ParseToken(tokenString string) (claims *models.Claims, err error) {
// 	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	claims, ok := token.Claims.(*models.Claims)

// 	if !ok {
// 		return nil, err
// 	}

// 	return claims, nil
// }

func ParseToken(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				return nil, errors.New("signature is invalid")
			}
		}
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
