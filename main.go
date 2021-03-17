package main

import (
	"ecommerce/views"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// views
var homeView *views.View
var contactView *views.View

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// load template
	err := homeView.Template.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// load template
	err := contactView.Template.Execute(w, nil)
	if err != nil {
		panic(err)
	}
}

func error404(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "Custom 404")
}

func main() {
	homeView = views.NewView("views/home.gohtml")
	contactView = views.NewView("views/contact.gohtml")

	// custom 404 handler
	var handler404 http.Handler = http.HandlerFunc(error404)

	r := mux.NewRouter()

	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)

	r.NotFoundHandler = handler404
	http.ListenAndServe(":3000", r)
}
