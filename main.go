package main

import (
	"log"
	"net/http"
	"os"
	"text/template"

	"forum/handler"
	"forum/model"
	"forum/service"
)

var templates = template.Must(template.ParseFiles("templates/signup.html", "templates/login.html"))
var globalUser *model.User = &model.User{}

func SignUp(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(globalUser.Name)

	switch r.Method {
	case http.MethodPost:
		user, err := service.CreateUser(
			r.FormValue("username"),
			r.FormValue("email"),
			r.FormValue("password"),
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = service.SignUp(user)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		err := templates.ExecuteTemplate(w, "signup.html", globalUser)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handler.MainPage)
	mux.HandleFunc("/signup/", SignUp)
	// mux.HandleFunc("/register/", Register)
	//http.FileServer(http.Dir("static"))

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
