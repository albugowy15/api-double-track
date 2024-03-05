package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func GetBody(w http.ResponseWriter, r *http.Request, v any) {
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		log.Printf("error decode body: %v", err)
		SendError(w, "bad request", http.StatusBadRequest)
		return
	}
}

type ErrorJsonResponse struct {
	Error string `json:"error"`
}

func SendError(w http.ResponseWriter, message string, status int) {
	res := ErrorJsonResponse{
		Error: message,
	}
	json, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("error marshaling json:%v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)
}

func SendJson(w http.ResponseWriter, res any, status int) {
	json, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("error marshaling json:%v", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(json)
}
