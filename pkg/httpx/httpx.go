package httpx

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var (
	ErrInternalServer = errors.New("internal server errror")
	ErrDecodeJsonBody = errors.New("error decode json body")
)

func GetBody(r *http.Request, v any) error {
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		log.Printf("error decode json body: %v", err)
		return ErrDecodeJsonBody
	}
	return nil
}

type ErrorJsonResponse struct {
	Error string `json:"error"`
}

type DataJsonResponse struct {
	Data interface{} `json:"data"`
}

type MessageJsonResponse struct {
	Message string `json:"message"`
}

func Send(w http.ResponseWriter, res any, status int) {
	json, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)
}

func SendError(w http.ResponseWriter, err error, status int) {
	res := ErrorJsonResponse{
		Error: err.Error(),
	}
	Send(w, res, status)
}

func SendMessage(w http.ResponseWriter, message string, status int) {
	res := MessageJsonResponse{
		Message: message,
	}
	Send(w, res, status)
}

func SendData(w http.ResponseWriter, data interface{}, status int) {
	res := DataJsonResponse{
		Data: data,
	}
	Send(w, res, status)
}
