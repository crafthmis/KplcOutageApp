package controllers

import (
	"kplc-outage-app/db"
	"kplc-outage-app/models"
	"kplc-outage-app/utils"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// The string "my_secret_key" is just an example and should be replaced with a secret key of sufficient length and complexity in a real-world scenario.
var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func Login(c *gin.Context) {

	var user models.Contact

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.Contact

	db.GetDB().Where("msisdn = ?", user.Msisdn).First(&existingUser)

	if existingUser.CntID == 0 {
		c.JSON(400, gin.H{"error": "user does not exist"})
		return
	}

	errHash := utils.CompareHashPassword(user.Password, existingUser.Password)

	if !errHash {
		c.JSON(400, gin.H{"error": "invalid password"})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &models.Claims{
		Role: existingUser.Role.Name,
		StandardClaims: jwt.StandardClaims{
			Subject:   existingUser.Msisdn,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		c.JSON(500, gin.H{"error": "could not generate token"})
		return
	}

	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)
	c.JSON(200, gin.H{"success": "user logged in"})
}

// PATH: go-auth/controllers/auth.go

func Signup(c *gin.Context) {
	var user models.Contact

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.Contact

	db.GetDB().Where("msisdn = ?", user.Email).First(&existingUser)

	if existingUser.SubID == 0 {
		c.JSON(400, gin.H{"error": "user already exists"})
		return
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(user.Password)

	if errHash != nil {
		c.JSON(500, gin.H{"error": "could not generate password hash"})
		return
	}

	db.GetDB().Create(&user)

	c.JSON(200, gin.H{"success": "user created"})
}
