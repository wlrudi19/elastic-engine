package jwt

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWT interface {
	GenerateAccessToken(userId int, email string) (string, error)
	ValidateToken(tokenString string) error
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
		"exp":   time.Now().Add(time.Minute * 1).Unix(), //time expired 1 menit
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

func (j *jwtImpl) ValidateToken(tokenString string) error {
	log.Printf("[JWT] validate tokenString, %s", tokenString)

	secretKey := []byte("x-elastic-engine")

	//parsing & validasi metode
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims) //return is bool

	if !ok || !token.Valid {
		return errors.New("token invalid")
	}

	// get payload map
	userId := int(claims["Id"].(float64))
	username := claims["Email"].(string)

	log.Println("user id:", userId)
	log.Println("email:", username)

	expTime := time.Unix(int64(claims["exp"].(float64)), 0)

	if expTime.Before(time.Now()) {
		return errors.New("token has expired")
	}

	return nil
}
