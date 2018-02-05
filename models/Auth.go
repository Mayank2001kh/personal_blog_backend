package models

import (
	"fmt"
	"net/http"
	"time"
	//"crypto/sha256"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type User struct {
	Id        int       `xorm:pk"`
	FirstName string    `xorm:"varchar(25) not null 'first_name'"`
	LastName  string    `xorm:"varchar(25) not null 'last_name'"`
	Email     string    `xorm:"varchar(25) not null unique 'usr_name'"`
	Password  string    `xorm:" varchar(25) not null 'password'"`
	Created   time.Time `xorm:"created"`
}

func (U User) Create(r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	form := r.Form
	fmt.Println(form)

	engine, _ := xorm.NewEngine("mysql", "root:admin@/golang_be?charset=utf8")

	fmt.Println(engine.Query("show tables;"))
	//new_User = User{}
}

func init() {
	engine, _ := xorm.NewEngine("mysql", "root:admin@/golang_be?charset=utf8")
	err := engine.Sync2(new(User))

	fmt.Println(err)

}
