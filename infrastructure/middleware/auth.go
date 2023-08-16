package middleware

// import (
// 	"log"
// 	"net/http"

// 	"github.com/golang-jwt/jwt"
// 	httputils "github.com/wlrudi19/elastic-engine/helper/http"
// )

// type Auth interface {
// 	ServeHTTP(writer http.ResponseWriter, request *http.Request)
// 	ValidateToken(tokenString string) (*jwt.Token, error)
// }

// type auth struct {
// 	Handler http.Handler
// }

// func NewAuth(handler http.Handler) Auth {
// 	return &auth{
// 		Handler: handler,
// 	}
// }

// func (au *auth) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
// 	tokenString := request.Header.Get("Authorization")
// 	if tokenString == "" {
// 		respon := []httputils.StandardError{
// 			{
// 				Code:   "401",
// 				Title:  "Unauthorized",
// 				Detail: "You are not authorized to access this resource",
// 				Object: httputils.ErrorObject{},
// 			},
// 		}
// 		httputils.WriteErrorResponse(writer, http.StatusBadRequest, respon)
// 		return
// 	}

// 	token, err := au.ValidateToken(tokenString)
// 	if err != nil || !token.Valid {
// 		respon := []httputils.StandardError{
// 			{
// 				Code:   "401",
// 				Title:  "Unauthorized",
// 				Detail: "You are not authorized to access this resource",
// 				Object: httputils.ErrorObject{},
// 			},
// 		}
// 		httputils.WriteErrorResponse(writer, http.StatusBadRequest, respon)
// 		return
// 	}

// 	au.Handler.ServeHTTP(writer, request)
// }

// func (au *auth) ValidateToken(tokenString string) (*jwt.Token, error) {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return []byte("x-elastic-engine"), nil
// 	})

// 	if err != nil {
// 		log.Printf("[JWT] failed to parse token, %v", err)
// 		return nil, err
// 	}

// 	return token, nil
// }
