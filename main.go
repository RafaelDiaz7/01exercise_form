package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var tpl *template.Template

type user struct {
	userName  string
	uEmail    string
	uPassword string
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func newUser(username string, email string, password string) *user {
	nUser := user{userName: username, uEmail: email, uPassword: password}
	return &nUser
}

func home(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "default.gohtml", nil)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)

	mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		//NOTE : Invoke ParseForm or ParseMultipartForm before reading form values
		r.ParseForm()
		/*
			Reads individual key-value pairs from
			r.Form object. Note that these include both data sent
			through request url and request body
		*/
		//newUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"))
		fmt.Printf("USERNAME => %s\n", newUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("password")).userName)
		tpl.ExecuteTemplate(w, "default.gohtml", newUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("password")))

		//var message = "All good!"
		//tpl.ExecuteTemplate(w, "default.gohtml", message)
	})

	http.ListenAndServe(":8080", mux)
}
