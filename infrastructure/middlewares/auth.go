package middlewares

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	httputils "github.com/wlrudi19/elastic-engine/helper/http"
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

		err := ValidateToken(tokenString)

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

		next.ServeHTTP(writer, request)
	})
}

func GenerateAccessToken(userId int, email string) (string, error) {
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

func ValidateToken(tokenString string) error {
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
	// harusnya buat validasi emailnya tapi ada masalah pas manggil fungsi FindUserLogic
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
