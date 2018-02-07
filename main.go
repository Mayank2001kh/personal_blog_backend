package main

import (
	"fmt"
	"log"
	"net/http"
	"personal_blog_backend/controllers"
)

// go get github.com/go-xorm/xorm
// go get github.com/go-sql-driver/mysql
// go get github.com/gorilla/sessions

func main() {

	fmt.Println("Project starts...")

	fs := http.FileServer(http.Dir("./media"))
	http.Handle("/media/", http.StripPrefix("/media/", fs))
	http.HandleFunc("/auth/register/", controllers.Register) // setting router rule
	http.HandleFunc("/auth/drop/", controllers.DropUser)
	http.HandleFunc("/auth/login/", controllers.Login)
	http.HandleFunc("/auth/logout/", controllers.Logout)
	http.HandleFunc("/auth/profile/", controllers.Profile)
	err := http.ListenAndServe(":9090", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
