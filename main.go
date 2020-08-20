package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Buena practica: agrupar todas las declaraciones de variables.
var (
	User *user
	tpl *template.Template
)

type user struct {
	// Cada campo debe ser exportado para que pueda ser accedido desde los
	// templates.
	UserName  string
	UEmail    string
	uPassword string
}

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func newUser(username string, email string, password string) *user {
	/*nUser := user{userName: username, uEmail: email, uPassword: password}
	return &nUser*/
	// de esta manera es mas legible y conciso:
	return &user{
		UserName:  username,
		UEmail:    email,
		uPassword: password,
	}
}

// Este handler es muy corto, deberia ser un closure y el handler de signup
// por el contrario, deberia ser una funcion independiente.
func home(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	// Consejo: utiliza librerias como gorilla/mux para evitar tener que analizar
	// el metodo HTTP con el cual fue realizada la peticion. En su lugar, harias
	// algo asi como router.HandleFunc("/dashboard", dashboard).Methods("GET")
	if r.Method != http.MethodGet {
		http.Error(w, "ERR_METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
		return
	}
	// La siguiente verificacion puede ser realizada por un "middleware". Te
	// recomiendo investigar sobre ello; se usa mucho en sistemas complejos.
	if User == nil {
		// User not logged in; redirect to home to signup.
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if (User.UEmail == "") || (User.UserName == "") || (User.uPassword == "") {
		// User not logged in; redirect to home to signup.
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "dashboard.gohtml", User)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", home)

	mux.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "ERR_METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed)
			return
		}
		//NOTE : Invoke ParseForm or ParseMultipartForm before reading form values
		// r.ParseForm es innecesario, solo es util cuando quieres analizar el
		// error que retorna, cosa que no se hace aqui.
		// r.ParseForm()
		/*
			Reads individual key-value pairs from
			r.Form object. Note that these include both data sent
			through request url and request body
		*/
		//newUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"))
		// Consejo: esto es muy dificil de leer, es mucho mejor asignar los campos
		// del formulario individualmente en su respectiva variable.
		// fmt.Printf("USERNAME => %s\n", newUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("password")).userName)
		// tpl.ExecuteTemplate(w, "default.gohtml", newUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("password")))
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Nota: User es una variable global -> es accesible por todas las
		// funciones. En el mundo real se utilizan cookies o JWT para guardar
		// datos de usuario y tener acceso a ellos entre handlers.
		User = newUser(username, email, password)
		fmt.Println("User:", *User)
		//var message = "All good!"
		//tpl.ExecuteTemplate(w, "default.gohtml", message)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	})

	mux.HandleFunc("/dashboard", dashboard)

	http.ListenAndServe(":8080", mux)
}
