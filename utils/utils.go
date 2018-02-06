package utils

import "net/http"
import "fmt"
import "personal_blog_backend/models"
import "encoding/json"

func AcceptMethod(r *http.Request, m string) bool {
	switch method := r.Method; method {
	case m:
		fmt.Println("Accepted")
		return true
	default:
		return false

		
	}
}