// PATH: go-auth/utils/GenerateHashPassword.go

package utils

import (
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
