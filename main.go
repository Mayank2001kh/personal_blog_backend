package main

import (
	"fmt"
	"log"
	"net/http"
	"personal_blog_backend/controllers"

	"github.com/gorilla/context"
	"github.com/rs/cors"
)

// go get github.com/go-xorm/xorm
// go get github.com/go-sql-driver/mysql
// go get github.com/gorilla/sessions
// go get github.com/gorilla/context
// go get github.com/rs.cors

func main() {

	fmt.Println("Project starts...")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PATCH"},
	})

	fs := http.FileServer(http.Dir("./media"))
	http.Handle("/media/", http.StripPrefix("/media/", fs))
	http.HandleFunc("/auth/register/", controllers.Register) // setting router rule
	http.HandleFunc("/auth/drop/", controllers.DropUser)
	http.HandleFunc("/auth/login/", controllers.Login)
	http.HandleFunc("/auth/logout/", controllers.Logout)
	http.HandleFunc("/auth/profile/", controllers.Profile)
	err := http.ListenAndServe(":9090", c.Handler(context.ClearHandler(http.DefaultServeMux))) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
