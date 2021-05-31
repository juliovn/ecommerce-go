package main

import (
	"ecommerce/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	// static pages handler
	staticController := controllers.NewStatic()

	// users handler
	usersController := controllers.NewUsers()

	r := mux.NewRouter()

	// HOME
	r.Handle("/", staticController.HomeView).Methods("GET")

	// CONTACT
	r.Handle("/contact", staticController.ContactView).Methods("GET")

	// SIGNUP
	r.HandleFunc("/signup", usersController.New).Methods("GET")
	r.HandleFunc("/signup", usersController.Create).Methods("POST")


	r.NotFoundHandler = staticController.Error404View
	http.ListenAndServe(":3000", r)
}
