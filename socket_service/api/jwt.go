package api

import (
	"fmt"
	"log"
	"net/http"

	service "github.com/Mahamed-Belkheir/scalechat-backend/socket_service"
	jwt "github.com/dgrijalva/jwt-go"
)

type claim struct {
	Username string
	jwt.StandardClaims
}

type JWT struct {
	secret string
}

func NewJWT(config service.Config) JWT {
	return JWT{
		secret: config.Secret,
	}
}

func (j JWT) verify(req *http.Request) (string, error) {
	tokenString := req.Header.Get("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("no authorization header provided")
	}
	tokenString = tokenString[7:]
	token, err := jwt.ParseWithClaims(tokenString, &claim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return "", err
	}
	claims := token.Claims.(*claim)
	log.Printf("debug: user connected with claims: %v", claims)
	return claims.Subject, nil
}
