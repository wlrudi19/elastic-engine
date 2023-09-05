package middlewares

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	httputils "github.com/wlrudi19/elastic-engine/helper/http"
)

type contextKey string

const (
	ContextKeyUserEmail contextKey = "userEmail"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		tokenString := request.Header.Get("Authorization")
		log.Printf("[MW] token string: %s", tokenString)

		if tokenString == "" {
			respon := []httputils.StandardError{
				{
					Code:   "401",
					Title:  "Unauthorized",
					Detail: "You are not authorized to access this resource",
					Object: httputils.ErrorObject{},
				},
			}
			httputils.WriteErrorResponse(writer, http.StatusBadRequest, respon)
			return
		}

		claims, err := ValidateToken(tokenString)

		if err != nil {
			respon := []httputils.StandardError{
				{
					Code:   "401",
					Title:  "Unauthorized",
					Detail: "Your access token invalid",
					Object: httputils.ErrorObject{},
				},
			}
			httputils.WriteErrorResponse(writer, http.StatusBadRequest, respon)
			return
		}

		email, ok := claims["Email"].(string)

		if !ok {
			respon := []httputils.StandardError{
				{
					Code:   "401",
					Title:  "Unauthorized",
					Detail: "Your access token invalid",
					Object: httputils.ErrorObject{},
				},
			}
			httputils.WriteErrorResponse(writer, http.StatusBadRequest, respon)
			return
		}

		ctx := context.WithValue(request.Context(), ContextKeyUserEmail, email)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GenerateAccessToken(userId int, email string) (string, error) {
	log.Printf("[JWT] generate access token with email: %s", email)
	//generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Id":    userId,
		"Email": email,
		"exp":   time.Now().Add(time.Minute * 10).Unix(), //time expired 1 menit
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

func ValidateToken(tokenString string) (map[string]interface{}, error) {
	log.Printf("[JWT] validate tokenString, %s", tokenString)

	secretKey := []byte("x-elastic-engine")

	//parsing & validasi metode hashing
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims) //return is bool

	if !ok || !token.Valid {
		return nil, errors.New("token invalid")
	}

	expTime := time.Unix(int64(claims["exp"].(float64)), 0)

	if expTime.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}
