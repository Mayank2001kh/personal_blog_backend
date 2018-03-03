package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Photo struct {
	Id          int64     `xorm:"pk autoincr"`
	UserId      int64     `xorm:"index"`
	Title       string    `xorm:"varchar(100) notnull 'title'"`
	Description string    `xorm:"varchar(400) notnull 'description'"`
	File        string    `xorm:"varchar(400) notnull 'file'"`
	Uploaded    time.Time `xorm:"created"`
}

type AbstractPhotoWrapper struct {
	AbstractUser  User  `xorm:"extends"`
	AbstractPhoto Photo `xorm:"extends"`
}

func (AbstractPhotoWrapper) TableName() string {
	return "user"
}

func (p Photo) Upload(r *http.Request, userid int64) (Photo, error) {

	// media/gallery/userid/fileid_filename
	engine, _ := xorm.NewEngine("mysql", connect_str)
	r.ParseMultipartForm(32 << 20)
	form := r.Form
	titleList := form["title"]
	descriptionList := form["description"]
	file, handler, err := r.FormFile("file")
	//verify file here
	if err != nil {
		return p, err
	}
	// verify form here
	if len(titleList)*len(descriptionList) <= 0 {
		return p, errors.New("Invalid form")
	}
	// Change the structure here

	p.Title = titleList[0]
	p.Description = descriptionList[0]
	p.File = "Pending"
	p.UserId = userid
	_, err = engine.Insert(&p)
	if err != nil {
		return p, err
	}
	file_url := fmt.Sprintf("/media/Gallery/%v/%v_%v", userid, p.Id, handler.Filename)
	p.File = file_url
	_, err = engine.Id(p.Id).Update(&p)
	if err != nil {
		return p, err
	}
	os.Mkdir("./media/Gallery", 0777)
	os.Mkdir(fmt.Sprintf("./media/Gallery/%v", userid), 0777)

	f, err := os.OpenFile("."+file_url, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return p, err
	}
	io.Copy(f, file)
	f.Close()
	file.Close()

	fmt.Println(titleList, descriptionList, file, handler, err)
	return p, nil
}

func (p Photo) Modify(r *http.Request) (Photo, string, error) {
	return p, "success", nil
}

func (p Photo) Delete(r *http.Request) (string, error) {
	engine, _ := xorm.NewEngine("mysql", connect_str)
	r.ParseMultipartForm(32 << 20)
	form := r.Form
	id_list := form["photoid"]
	if len(id_list) <= 0 {
		return "Invalid Form", errors.New("Invalid Form")
	} else {
		pid := id_list[0]
		// photo id

		engine.Id(pid).Get(&p)
		// query photo from database
		fmt.Println(p)
		file_dir := p.File
		os.Remove("." + file_dir)
		_, err := engine.Id(pid).Delete(p)
		if err != nil {
			return "Database Error", err
		} else {
			return "Success", nil
		}
	}
	//
	return "None", nil
}

func (p Photo) Fetch(r *http.Request, userid int64) ([]Photo, error) {
	// Fetch a photo (or photos) and return a slice of Photo along with error
	engine, _ := xorm.NewEngine("mysql", connect_str)
	r.ParseMultipartForm(32 << 20)
	form := r.Form
	start_list := form["start"]
	number_list := form["number"]
	type_list := form["type"]
	id_list := form["photoid"]
	photo_slice := []Photo{Photo{}}
	if len(start_list)*len(number_list)*len(type_list) <= 0 {
		return photo_slice, errors.New("Invalid Form")
	} else {
		switch type_list[0] {
		case "single":
			// need to provide photo id
			_, err := engine.Id(id_list[0]).Get(&p)
			if err != nil {
				return []Photo{p}, err
			} else {
				return []Photo{p}, nil
			}

		case "multiple":
			// need to provide start, number
			wrapper := make([]Photo, 0)
			start, _ := strconv.ParseInt(start_list[0], 10, 32)
			number, _ := strconv.ParseInt(number_list[0], 10, 32)

			err := engine.Limit(int(number), int(start)).Find(&wrapper)
			if err != nil {
				return wrapper, err
			} else {
				return wrapper, nil
			}

		case "self":
			wrapper := make([]AbstractPhotoWrapper, 0)
			start, _ := strconv.ParseInt(start_list[0], 10, 32)
			number, _ := strconv.ParseInt(number_list[0], 10, 32)
			err := engine.Join("INNER", "photo", fmt.Sprintf("photo.user_id = %v", userid)).Limit(int(number), int(start)).Find(&wrapper)
			plist := make([]Photo, len(wrapper))
			for i, _ := range plist {
				plist[i] = wrapper[i].AbstractPhoto
			}
			if err != nil {
				return plist, err
			} else {
				return plist, nil
			}
			// need to provide user id, start, number

		default:
			return []Photo{Photo{}}, errors.New("Type not supported")

		}

	}

	photoList := []Photo{}
	return photoList, nil
}

func init() {
	engine, _ := xorm.NewEngine("mysql", connect_str)
	err := engine.Sync2(new(Photo))
	if err != nil {
		fmt.Println(err)
	}
	//_, err = engine.Exec("ALTER TABLE photo DROP FOREIGN KEY IDX_photo_user_id;")
	// _, err = engine.Exec("ALTER TABLE photo ADD FOREIGN KEY IDX_photo_user_id(user_id) REFERENCES user(id) ON DELETE CASCADE;")
	// fmt.Println(err)

	fmt.Println("Gallery loaded")
}
