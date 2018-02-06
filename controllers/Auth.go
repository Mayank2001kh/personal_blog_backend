package controllers

import (
	"encoding/json"
	"fmt"
	"personal_blog_backend/models"
	"net/http"
	"github.com/gorilla/sessions"
	"personal_blog_backend/config"
	
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func Register(w http.ResponseWriter, r *http.Request) {

	response := models.BaseResponse{"success", "Hello"}
	enc := json.NewEncoder(w)
	switch method := r.Method; method {
		case "POST": 
			newUser := models.User{}
			message,err := newUser.Create(r)
			if err != nil {
				response = models.BaseResponse{"error", message}
				enc.Encode(response)
			} else {
				enc.Encode(response)
			}
		default: 
			response = models.BaseResponse{"error", fmt.Sprintf("Method: %v not supported", method)}
			enc.Encode(response)
	}
	

	
}

func DropUser(w http.ResponseWriter, r *http.Request) {

	response := models.BaseResponse{"success", "User has been dropped"}
	enc := json.NewEncoder(w)
	switch method := r.Method; method {
	case "POST":
		newUser := new(models.User)
		message, err :=newUser.DropUser(r)
		if err != nil {
			response = models.BaseResponse{"error", message}
		}
		enc.Encode(response)
		
	
	default:
		response = models.BaseResponse{"error", fmt.Sprintf("Method: %v not supported", method)}
		enc.Encode(response)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	
	response := models.BaseResponse{"success", "Login successfully"}
	enc := json.NewEncoder(w)
	switch method := r.Method; method {
	case "POST":
		newUser := new(models.User)
		uid,err := newUser.Authenticate(r)
		if err != nil {
			response = models.BaseResponse{"error","Authentication fail"}
		} else {
			session, _ := store.Get(r, configs.COOKIENAME)
			if session.Values["uid"] == uid {
				response = models.BaseResponse{"error","Do not log in repeatedly"}
			}
			session.Values["uid"] = uid
			session.Save(r,w)
			

		}

		enc.Encode(response)
	default:
		response = models.BaseResponse{"error", fmt.Sprintf("Method: %v not supported", method)}
		enc.Encode(response)


	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	
	response := models.BaseResponse{"success", "Logout successfully"}
	enc := json.NewEncoder(w)

	switch method := r.Method; method {
	case "GET":
		session, _ := store.Get(r, configs.COOKIENAME)
		uid := session.Values["uid"]
		if uid == nil {
			response = models.BaseResponse{"error", "You are not logged in"}
		} else {
			session.Values["uid"] = nil
			err := session.Save(r,w)
			if err != nil {
				response = models.BaseResponse{"error", fmt.Sprint(err)}
			}
		}
			enc.Encode(response)
		
	default:
		response = models.BaseResponse{"error", fmt.Sprintf("Method: %v not supported", method)}
		enc.Encode(response)
	}
	
}
