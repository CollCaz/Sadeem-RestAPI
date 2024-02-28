package auth

import (
	"Sadeem-RestAPI/internal/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	Name   string `json:"name"`
	UserID int    `json:"-"`
	Admin  bool   `json:"admin"`
	jwt.RegisteredClaims
}

func CreateJwtToken(user *models.User) (string, error) {
	signingKey := os.Getenv("JWT_SIGNING_KEY")
	claims := &JwtClaims{
		user.UserName,
		user.ID,
		user.IsAdmin,
		jwt.RegisteredClaims{
			ID:        "user_token",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return rawToken.SignedString([]byte(signingKey))
}
