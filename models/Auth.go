package models

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"personal_blog_backend/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var connect_str = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", configs.DBUSER, configs.DBPASS, configs.DBHOST, configs.DBPORT, configs.DBNAME)

func PassHash(password []byte) string {

	h := sha256.New()
	h.Write(password)
	pass_hashed := fmt.Sprintf("%x", h.Sum(nil))

	return pass_hashed
}

func Validate(form map[string][]string) (string, string, string, string, string, error) {
	email := form["email"]
	firstname := form["firstname"]
	lastname := form["lastname"]
	password := form["password"]
	username := form["username"]
	if len(email)*len(firstname)*len(lastname)*len(password)*len(username) == 0 {
		return "", "", "", "", "", errors.New("Invalid form")

	} else {
		return username[0], password[0], email[0], firstname[0], lastname[0], nil
	}

}

type User struct {
	Id        int64     `xorm:"pk autoincr"`
	UserName  string    `xorm:"varchar(25) notnull unique 'username'"`
	FirstName string    `xorm:"varchar(25) notnull 'firstname'"`
	LastName  string    `xorm:"varchar(25) notnull 'lastname'"`
	Email     string    `xorm:"varchar(40) notnull unique 'email'"`
	Password  string    `xorm:"varchar(256) notnull 'password'"`
	Created   time.Time `xorm:"created"`
}

// show columns from user

type Group struct {
}

type UserProfile struct {
	Id     int64  `xorm:"pk autoincr"`
	UserId int64  `xorm:"index"`
	Avatar string `xorm:"varchar(256) 'avatar'"`
	Intro  string `xorm:"text 'intro'"`
}

type AbstractUserWrapper struct {
	AbstractUser    User        `xorm:"extends"`
	AbstractProfile UserProfile `xorm:"extends"`
}

func (AbstractUserWrapper) TableName() string {
	return "user"
}

func (U User) Create(r *http.Request) (string, error) {
	// Step1: parse form
	// Step2: validate form
	// Step3: hash the password
	// Step4: Create the user instance and dump to database
	// Step5: Create profile for the user
	// Step6: Return message and error (if applicable)

	r.ParseMultipartForm(32 << 20)
	form := r.Form

	engine, _ := xorm.NewEngine("mysql", connect_str)

	username, password, email, firstname, lastname, err := Validate(form)
	if err != nil {
		return "Invalid form", err
	} else {

		pass_hashed := PassHash([]byte(password))

		new_user := User{UserName: username,
			FirstName: firstname,
			LastName:  lastname,
			Email:     email,
			Password:  pass_hashed}
		fmt.Println(new_user)
		affected, err := engine.Insert(&new_user)
		if err != nil {

			fmt.Println(err)
			return fmt.Sprint(err), err
			// insert entry and return possible error

		} else {
			// successfully inserted user info, now create the profile
			new_profile := UserProfile{UserId: new_user.Id,
				Avatar: "media/def.jpg",
				Intro:  "None"}
			affected, err = engine.Insert(&new_profile)

			fmt.Println(affected, new_user.Id)
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

func (U User) Authenticate(r *http.Request) (int64, error) {
	// Step1: Parse form
	// Step2: Check if username and password matches
	// Step3: Return userid and error
	engine, _ := xorm.NewEngine("mysql", connect_str)
	r.ParseMultipartForm(32 << 20)
	form := r.Form
	fmt.Println(form)
	username := form["username"]
	password := form["password"]
	if len(username)*len(password) <= 0 {
		fmt.Println(len(username), len(password))
		return 0, errors.New("Invalid form")
	}

	hashed_password := PassHash([]byte(password[0]))
	user := &User{UserName: username[0]}
	has, _ := engine.Get(user)

	if has == false {
		return 0, errors.New("User does not exist")
	} else if user.Password == hashed_password {
		return user.Id, nil

	} else {
		return 0, errors.New("username/password unmatch")
	}

	return 0, errors.New("unexpected error")
}

func (U User) FetchProfile(userId int64) (UserProfile, error) {
	// Step1: Get userid
	// Step2: Join query
	// Step3: Return the profile
	engine, _ := xorm.NewEngine("mysql", connect_str)
	wrapper := make([]AbstractUserWrapper, 0)

	err := engine.Join("INNER", "user_profile", "user_profile.user_id = user.id ").Find(&wrapper)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(wrapper)
	}

	return wrapper[0].AbstractProfile, err
}

func (U User) SetProfile(r *http.Request, userid int64) (UserProfile, error) {
	// Set avatar and intro of the profile
	// Step1: Parse form and extract file
	// Step2: Set profile photo by copying imgage and set the path
	// Step3: Set intro
	// Step4: Return the user profile
	r.ParseMultipartForm(32 << 20)
	form := r.Form
	intro := "None"
	if len(form["intro"]) <= 0 {
		fmt.Println(form)

	} else {
		intro = form["intro"][0]
	}

	file, handler, err := r.FormFile("file")
	fmt.Println(err)
	fmt.Println(r.MultipartForm)
	if err != nil {

		return UserProfile{}, errors.New(fmt.Sprint(err))
	}

	avatar := "/media/avatars/user_" + fmt.Sprint(userid) + "_avatar_" + handler.Filename
	f, err := os.OpenFile("."+avatar, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return UserProfile{}, errors.New(fmt.Sprint(err))

	}
	io.Copy(f, file)
	fmt.Println(err)
	f.Close()
	file.Close()

	engine, _ := xorm.NewEngine("mysql", connect_str)
	wrapper := make([]AbstractUserWrapper, 0)
	err = engine.Join("INNER", "user_profile", "user_profile.user_id = user.id ").Find(&wrapper)

	if err != nil {
		fmt.Println(err)
		return UserProfile{}, err
	}

	profile := wrapper[0].AbstractProfile
	profile.Avatar = avatar
	profile.Intro = intro
	affected, err := engine.Id(profile.Id).Update(&profile)
	if err != nil {
		return UserProfile{}, err
	}
	fmt.Println("Updated, affected rows: ", affected)

	return profile, err

}

func init() {
	// connect to mysql
	//connect_str := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8",configs.DBUSER,configs.DBPASS,configs.DBHOST,configs.DBPORT,configs.DBNAME)
	engine, _ := xorm.NewEngine("mysql", connect_str)
	err := engine.Sync2(new(User))
	err = engine.Sync2(new(UserProfile))
	_, err = engine.Exec("ALTER TABLE user_profile ADD FOREIGN KEY " +
		"IDX_user_profile_user_id(user_id) REFERENCES user(id) ON DELETE CASCADE;")
	fmt.Println(err)
}
