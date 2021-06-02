package main

import (
	"ecommerce/controllers"
	"ecommerce/models"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	host	 = "localhost"
	port	 = 5432
	user	 = "postgres"
	password = "postgres"
	dbname	 = "ecommerce_dev"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	// DB connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.AutoMigrate()

	// uncomment if DB wipe is needed
	//us.DestructiveReset()

	// static pages handler
	staticController := controllers.NewStatic()

	// users handler
	usersController := controllers.NewUsers(us)

	r := mux.NewRouter()

	// HOME
	r.Handle("/", staticController.HomeView).Methods("GET")

	// CONTACT
	r.Handle("/contact", staticController.ContactView).Methods("GET")

	// SIGNUP
	r.HandleFunc("/signup", usersController.New).Methods("GET")
	r.HandleFunc("/signup", usersController.Create).Methods("POST")

	// LOGIN
	r.Handle("/login", usersController.LoginView).Methods("GET")
	r.HandleFunc("/login", usersController.Login).Methods("POST")

	// 404
	r.NotFoundHandler = staticController.Error404View

	// Server start
	http.ListenAndServe(":3000", r)
}
