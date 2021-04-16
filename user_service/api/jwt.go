package api

import (
	"fmt"
	"time"

	service "github.com/Mahamed-Belkheir/scalechat-backend/user_service"
	jwt "github.com/dgrijalva/jwt-go"
)

type JWT struct {
	secret string
}

func NewJWT(config service.Config) JWT {
	return JWT{
		secret: config.Secret,
	}
}

type claim struct {
	Username string
	jwt.StandardClaims
}

func (j JWT) sign(user *service.User) (string, error) {
	claims := claim{
		user.Username,
		jwt.StandardClaims{
			Subject:   user.ID,
			ExpiresAt: time.Now().Add(time.Hour * 4).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("error signing token %w", err)
	}
	return tokenString, nil
}
