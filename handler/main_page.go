package handler

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("templates/signup.html", "templates/login.html"))

func MainPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../templates/main.html")
}

// func SignUp(w http.ResponseWriter, r *http.Request) {
// 	// fmt.Println("$e")
// 	err := templates.ExecuteTemplate(w, "signup.html", nil)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }
