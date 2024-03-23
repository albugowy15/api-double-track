package httputil

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var ErrInternalServer = errors.New("internal server errror")

func GetBody(w http.ResponseWriter, r *http.Request, v any) {
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		log.Printf("error decode body: %v", err)
		SendError(w, errors.New("bad request"), http.StatusBadRequest)
		return
	}
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

func SendError(w http.ResponseWriter, err error, status int) {
	res := ErrorJsonResponse{
		Error: err.Error(),
	}
	json, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("error marshaling json:%v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)
}

func SendMessage(w http.ResponseWriter, message string, status int) {
	res := MessageJsonResponse{
		Message: message,
	}
	json, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("error marshaling json:%v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)
}

func SendData(w http.ResponseWriter, data interface{}, status int) {
	res := DataJsonResponse{
		Data: data,
	}
	json, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("error marshaling json:%v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)
}
