package middlewares

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	httputils "github.com/wlrudi19/elastic-engine/helper/http"
)

type Auth interface {
	Authenticate(http.Handler) http.Handler
	ValidateToken(tokenString string) (*jwt.Token, error)
}

type auth struct {
}

func NewAuth() Auth {
	return &auth{}
}

func (au *auth) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		//tokenString := request.Header.Get("Authorization")
		tokenCookie, err := request.Cookie("access-token")
		log.Printf("ini token %v", tokenCookie)
		if err != nil {
			log.Printf("ini error %v", err)
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

		//get token value
		tokenString := tokenCookie.Value
		log.Printf("ini token 2 %s", tokenCookie)
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

		token, err := au.ValidateToken(tokenString)
		log.Printf("ini token 3%v", token)
		if err != nil || !token.Valid {
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

		next.ServeHTTP(writer, request)
	})
}

func (au *auth) ValidateToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("x-elastic-engine"), nil
	})

	if err != nil {
		log.Printf("[JWT] failed to parse token, %v", err)
		return nil, err
	}

	return token, nil
}
