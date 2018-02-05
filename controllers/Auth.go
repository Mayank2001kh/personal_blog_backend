package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"my_be/models"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	response := models.BaseResponse{"success", "Hello"}
	enc := json.NewEncoder(w)

	newUser := models.User{}
	newUser.Create(r)

	enc.Encode(response)
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}
