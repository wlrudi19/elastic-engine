package jwt

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWT interface {
	GenerateAccessToken(userId int, email string) (string, error)
}

type jwtImpl struct {
}

func NewJWT() JWT {
	return &jwtImpl{}
}

func (j *jwtImpl) GenerateAccessToken(userId int, email string) (string, error) {
	log.Printf("[JWT] generate access token with email: %s", email)
	//generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Id":    userId,
		"Email": email,
		"exp":   time.Now().Add(time.Hour * 1).Unix(), //time expired 1 jam
	})

	//tandatangan token dengan secret key
	secretKey := []byte("x-elastic-engine")
	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		log.Printf("[JWT] failed to generate access token, %v", err)
		return "", err
	}

	return tokenString, nil
}
