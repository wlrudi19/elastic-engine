package httputils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ResponseDecorator interface {
	Decorate(w http.ResponseWriter)
}

type ContentTypeDecorator string

func (d ContentTypeDecorator) Decorate(w http.ResponseWriter) {
	w.Header().Set("Content-Type", string(d))
}

func NewContentTypeDecorator(contentType string) ContentTypeDecorator {
	return ContentTypeDecorator(contentType)
}

type CORSDecorator struct {
	allowedOrigin string
}

func (d *CORSDecorator) Decorate(w http.ResponseWriter) {
	//w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
}

func WriteResponse(w http.ResponseWriter, data []byte, status int, decorators ...ResponseDecorator) (int, error) {
	for _, decorator := range decorators {
		decorator.Decorate(w)
	}
	w.WriteHeader(status)
	return w.Write(data)
}

func WriteErrorResponse(w http.ResponseWriter, code int, errs []StandardError) {

	contentType := NewContentTypeDecorator("application/json")
	response := StandardEnvelope{
		Errors: errs,
	}
	errResponse, err := json.Marshal(response)
	if err != nil {
		WriteResponse(w, []byte(fmt.Sprintf(`{"errors":[{"code":"500","title":"Internal Server Error","detail":"%s","object":{"text":null,"type":0}}]}`, err.Error())), http.StatusInternalServerError, contentType)
	}

	WriteResponse(w, errResponse, code, contentType)
}
