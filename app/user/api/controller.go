package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/wlrudi19/elastic-engine/app/user/model"
	"github.com/wlrudi19/elastic-engine/app/user/service"
	httputils "github.com/wlrudi19/elastic-engine/helper/http"
)

type UserHandler interface {
	FindUserHandler(writer http.ResponseWriter, req *http.Request)
}

type userhandler struct {
	UserLogic service.UserLogic
}

func NewUserHandler(userLogic service.UserLogic) UserHandler {
	return &userhandler{
		UserLogic: userLogic,
	}
}

func (h *userhandler) FindUserHandler(writer http.ResponseWriter, req *http.Request) {
	var jsonReq model.UserRequest

	err := json.NewDecoder(req.Body).Decode(&jsonReq)

	if err != nil {
		respon := []httputils.StandardError{
			{
				Code:   "400",
				Title:  "Bad Request",
				Detail: "Permintaan tidak valid. Format JSON tidak sesuai",
				Object: httputils.ErrorObject{},
			},
		}
		httputils.WriteErrorResponse(writer, http.StatusBadRequest, respon)
		return
	}

	var user = model.UserResponse{}
	user, err = h.UserLogic.FindUserLogic(context.TODO(), jsonReq.Email)

	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			respon := []httputils.StandardError{
				{
					Code:   "404",
					Title:  "Not found",
					Detail: "User not found",
					Object: httputils.ErrorObject{},
				},
			}
			httputils.WriteErrorResponse(writer, http.StatusInternalServerError, respon)
			return
		}

		respon := []httputils.StandardError{
			{
				Code:   "500",
				Title:  "Internal server error",
				Detail: "Terjadi kesalahan internal pada server",
				Object: httputils.ErrorObject{},
			},
		}
		httputils.WriteErrorResponse(writer, http.StatusInternalServerError, respon)
		return
	}

	status := httputils.StandardStatus{
		ErrorCode: 200,
		Message:   "User finding successfully",
	}

	envelope := httputils.StandardEnvelope{
		Status: &status,
		Data:   &user,
	}

	responFix, err := json.Marshal(envelope)
	if err != nil {
		respon := []httputils.StandardError{
			{
				Code:   "500",
				Title:  "Internal server error",
				Detail: "Terjadi kesalahan internal pada server",
				Object: httputils.ErrorObject{},
			},
		}
		httputils.WriteErrorResponse(writer, http.StatusInternalServerError, respon)
		return
	}

	contentType := httputils.NewContentTypeDecorator("application/json")
	httpStatus := http.StatusCreated

	httputils.WriteResponse(writer, responFix, httpStatus, contentType)
}
