package models

import (
	"fmt"
	"net/http"
	"time"
	"crypto/sha256"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"personal_blog_backend/config"
)

var connect_str =  fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8",configs.DBUSER,configs.DBPASS,configs.DBHOST,configs.DBPORT,configs.DBNAME)

func PassHash(password []byte) string{
	h := sha256.New()
	h.Write(password)
	pass_hashed := fmt.Sprintf("%x", h.Sum(nil))

	return pass_hashed
}

func Validate(form map[string][]string) (string,string,string,string,string,error) {
	email := form["email"]
	firstname := form["firstname"]
	lastname := form["lastname"]
	password := form["password"]
	username := form["username"]
	if len(email)*len(firstname)*len(lastname)*len(password)*len(username) == 0 {
		return "","","","","",errors.New("Invalid form")

	} else {
		return username[0],password[0],email[0],firstname[0],lastname[0],nil
	}



}

type User struct {
	Id        int       `xorm:"pk autoincr"`
	UserName  string	`xorm:"varchar(25) notnull unique 'username'"`
	FirstName string    `xorm:"varchar(25) notnull 'firstname'"`
	LastName  string    `xorm:"varchar(25) notnull 'lastname'"`
	Email     string    `xorm:"varchar(40) notnull unique 'email'"`
	Password  string    `xorm:"varchar(256) notnull 'password'"`
	Created   time.Time `xorm:"created"`
}
// show columns from user

func (U User) Create(r *http.Request) (string, error) {
	// Step1: parse form
	// Step2: validate form by checking 
	// Step3: hash the password using sha256 algorithm
	// Step4: Create the user instance and dump to database
	// Step5: Return message and error

	r.ParseMultipartForm(32 << 20)
	form := r.Form
	//fmt.Println(form)

	engine, _ := xorm.NewEngine("mysql", connect_str)
	//fmt.Println(engine.Query("show tables;"))

	username,password,email,firstname,lastname,err := Validate(form) 
	if err != nil {
		return "Invalid form", err
	} else {
		
		pass_hashed := PassHash([]byte(password))
		//fmt.Println(pass_hashed,username,password,email,firstname,lastname,err)
		
		
		new_user := User{UserName:username,
			FirstName:firstname,
			LastName:lastname,
			Email:email,
			Password:pass_hashed}
		fmt.Println(new_user)
		affected, err := engine.Insert(&new_user)
		if err != nil {
			fmt.Println(err)
			return fmt.Sprint(err), err
			// insert entry and return possible error
			
		} else {
			fmt.Println(affected)
		}

	}
	

	return "success", nil
}

func (U User) DropUser(r *http.Request) (string, error) {
	// Step1: Parse form
	// Step2: Validate for necessary fields
	// Step3: Drop the user
	r.ParseMultipartForm(32 << 20)
	form := r.Form
	engine, _ := xorm.NewEngine("mysql", connect_str)
	user_id := form["userid"]
	if len(user_id) == 0 {
		return "Invalid form", errors.New("Invalid form")
	}

	user := new(User)
	affected, err := engine.Id(user_id[0]).Delete(user)

	if err != nil {
		return fmt.Sprint(affected), err
	}

	return "success", nil

}

func (U User) Authenticate(r *http.Request) (int, error) {
	// Step1: Parse form
	// Step2: Check if username and password matches
	// Step3: Return userid and error
	engine, _ := xorm.NewEngine("mysql", connect_str)
	r.ParseMultipartForm(32 << 20)
	form := r.Form
	username := form["username"]
	password := form["password"]
	if len(username) * len(password) <= 0 {
		return 0,errors.New("Invalid form")
	}
		
	hashed_password := PassHash([]byte(password[0]))
	user := &User{UserName:username[0]}
	has, _ := engine.Get(user)
	
	if has == false {
		return 0,errors.New("User does not exist")
	} else if user.Password == hashed_password {
		return user.Id, nil
		
	} else {
		return 0,errors.New("username/password unmatch")
	}

	
	return 0, errors.New("unexpected error")
}

func init() {
	// connect to mysql
	//connect_str := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8",configs.DBUSER,configs.DBPASS,configs.DBHOST,configs.DBPORT,configs.DBNAME)
	engine, _ := xorm.NewEngine("mysql", connect_str)
	err := engine.Sync2(new(User))

	fmt.Println(err)
}
