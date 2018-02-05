package main

import (
	"fmt"
	"log"
	"my_be/controllers"
	"net/http"
)

// go get github.com/go-xorm/xorm

func main() {

	fmt.Println("Project starts...")
	http.HandleFunc("/register/", controllers.Register) // setting router rule
	http.HandleFunc("/login/", controllers.Login)
	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
