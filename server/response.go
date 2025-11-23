package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, statusCode int, payload any) error {
	res, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	fmt.Println(string(res))
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(res)
	return nil
}

func respondWithErrorJSON(w http.ResponseWriter, statusCode int, err error) error {
	return respondWithJSON(w, statusCode, map[string]string{"error": err.Error()})
}
