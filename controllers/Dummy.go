package controllers

import (
	"encoding/json"
	"net/http"
	"personal_blog_backend/models"
)

func DummyAPI(w http.ResponseWriter, r *http.Request) {

	response := models.BaseResponse{}
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:3000")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	enc := json.NewEncoder(w)
	enc.Encode(response)
}
