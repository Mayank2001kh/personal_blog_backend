package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"personal_blog_backend/config"
	"personal_blog_backend/models"
)

type PhotoResponse struct {
	models.BaseResponse
	models.Photo
}

type PhotoListResponse struct {
	models.BaseResponse
	InnerData []models.Photo
}

func PhotoUpload(w http.ResponseWriter, r *http.Request) {
	p := models.Photo{}
	response := PhotoResponse{models.BaseResponse{Message: "Upload Success", Status: "success"}, p}
	enc := json.NewEncoder(w)

	switch method := r.Method; method {
	case "POST":
		session, _ := store.Get(r, configs.COOKIENAME)

		if session.Values["uid"] == nil {
			response = PhotoResponse{models.BaseResponse{Message: "you are not logged in", Status: "error"}, p}
		} else {
			photo, err := p.Upload(r, session.Values["uid"].(int64))
			if err != nil {
				response = PhotoResponse{models.BaseResponse{Message: fmt.Sprint(err), Status: "error"}, p}
			} else {
				response = PhotoResponse{models.BaseResponse{Message: fmt.Sprint(err), Status: "success"}, photo}
			}
			fmt.Println(p, photo)
		}
		enc.Encode(response)

	default:
		response := models.BaseResponse{"error", fmt.Sprintf("Method: %v not supported", method)}
		enc.Encode(response)
	}

}

func PhotoDrop(w http.ResponseWriter, r *http.Request) {
	response := models.BaseResponse{Status: "success", Message: "Drop photo successfully"}
	enc := json.NewEncoder(w)
	switch method := r.Method; method {
	case "POST":
		session, _ := store.Get(r, configs.COOKIENAME)
		if session.Values["uid"] == nil {
			response = models.BaseResponse{Message: "you are not logged in", Status: "error"}
		} else {
			p := models.Photo{}
			_, err := p.Delete(r)
			if err != nil {
				response = models.BaseResponse{Status: "error", Message: fmt.Sprint(err)}

			}
		}

		enc.Encode(response)

	default:
		response := models.BaseResponse{"error", fmt.Sprintf("Method: %v not supported", method)}
		enc.Encode(response)
	}

}

func PhotoFetch(w http.ResponseWriter, r *http.Request) {
	p := models.Photo{}

	enc := json.NewEncoder(w)
	response := PhotoListResponse{models.BaseResponse{Message: "Fetch Success", Status: "success"}, []models.Photo{p}}
	switch method := r.Method; method {
	case "POST":
		session, _ := store.Get(r, configs.COOKIENAME)
		if session.Values["uid"] == nil {
			response = PhotoListResponse{models.BaseResponse{Message: "You are not logged in", Status: "error"}, []models.Photo{p}}
		} else {
			plist, err := p.Fetch(r, session.Values["uid"].(int64))
			if err != nil {
				response = PhotoListResponse{models.BaseResponse{Message: fmt.Sprint(err), Status: "error"}, []models.Photo{p}}
			} else {
				response = PhotoListResponse{models.BaseResponse{Message: "Fetch Success", Status: "success"}, plist}
				fmt.Println(plist)
			}
		}

		enc.Encode(response)

	default:
		response := models.BaseResponse{"error", fmt.Sprintf("Method: %v not supported", method)}
		enc.Encode(response)
	}
}
