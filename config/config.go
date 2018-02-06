package configs

var HOST = "win"



var DBHOST string
var DBPORT string
var DBUSER string
var DBPASS string
var DBNAME string

var COOKIENAME = "johnnysbackend"

func init() {

	if HOST == "win" {
		 DBHOST = "localhost"
		 DBPORT = "5432"
		 DBUSER = "root"
		 DBPASS = "admin"
		 DBNAME = "golang_be"

	} else {

		 DBHOST = "localhost"
		 DBPORT = "3306"
		 DBUSER = "root"
		 DBPASS = "admin"
		 DBNAME = "golang_be"
	}
}