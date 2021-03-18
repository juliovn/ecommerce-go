package main

import (
	"ecommerce/views"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	// views
	homeView 	*views.View
	contactView *views.View
	signupView	*views.View
)



func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// load template
	must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// load template
	must(contactView.Render(w, nil))
}

func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(signupView.Render(w, nil))
}

func error404(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "Custom 404")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	homeView = views.NewView("base", "views/home.gohtml")
	contactView = views.NewView("base", "views/contact.gohtml")
	signupView = views.NewView("base", "views/signup.gohtml")

	// custom 404 handler
	var handler404 http.Handler = http.HandlerFunc(error404)

	r := mux.NewRouter()

	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/signup", signup)

	r.NotFoundHandler = handler404
	http.ListenAndServe(":3000", r)
}
