package middlewares

import (
	"log"
	"net/http"

	httputils "github.com/wlrudi19/elastic-engine/helper/http"
	"github.com/wlrudi19/elastic-engine/infrastructure/jwt"
)

type Auth interface {
	Authenticate(http.Handler) http.Handler
}

type auth struct {
}

func NewAuth() Auth {
	return &auth{}
}

func (au *auth) Authenticate(next http.Handler) http.Handler {
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

		err := jwt.NewJWT().ValidateToken(tokenString)

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
