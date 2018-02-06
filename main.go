package main

import (
	"fmt"
	"log"
	"personal_blog_backend/controllers"
	"net/http"
)

// go get github.com/go-xorm/xorm
// got get github.com/go-sql-driver/mysql

func main() {

	fmt.Println("Project starts...")
	http.HandleFunc("/auth/register/", controllers.Register) // setting router rule
	http.HandleFunc("/auth/drop/", controllers.DropUser)
	http.HandleFunc("/auth/login/", controllers.Login)
	http.HandleFunc("/auth/logout/", controllers.Logout)
	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
