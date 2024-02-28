package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

func CreateJwtToken(userName string, isAdmin bool) (string, error) {
	signingKey := os.Getenv("JWT_SIGNING_KEY")
	claims := &JwtClaims{
		userName,
		isAdmin,
		jwt.RegisteredClaims{
			ID:        "user_token",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return rawToken.SignedString([]byte(signingKey))
}
