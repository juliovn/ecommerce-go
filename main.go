package main

import (
	"ecommerce/controllers"
	"ecommerce/middleware"
	"ecommerce/models"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "ecommerce_dev"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	// DB connection
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	// uncomment if DB wipe is needed
	//services.DestructiveReset()

	r := mux.NewRouter()

	// static pages handler
	staticController := controllers.NewStatic()

	// users handler
	usersController := controllers.NewUsers(services.User)

	// items handler
	itemsController := controllers.NewItems(services.Item, r)

	// middleware
	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}

	// HOME
	r.Handle("/", staticController.HomeView).Methods("GET")

	// SIGNUP
	r.HandleFunc("/signup", usersController.New).Methods("GET")
	r.HandleFunc("/signup", usersController.Create).Methods("POST")

	// LOGIN
	r.Handle("/login", usersController.LoginView).Methods("GET")
	r.HandleFunc("/login", usersController.Login).Methods("POST")

	// ITEMS
	newItem := requireUserMw.Apply(itemsController.New)
	r.Handle("/items/new", newItem).Methods("GET")
	createItem := requireUserMw.ApplyFn(itemsController.Create)
	r.HandleFunc("/items", createItem).Methods("POST")
	showItem := requireUserMw.ApplyFn(itemsController.Show)
	r.HandleFunc("/items/{id:[0-9]+}", showItem).Methods("GET").Name(controllers.ShowItem)
	editItem := requireUserMw.ApplyFn(itemsController.Edit)
	r.HandleFunc("/items/{id:[0-9]+}/edit", editItem).Methods("GET")
	updateItem := requireUserMw.ApplyFn(itemsController.Update)
	r.HandleFunc("/items/{id:[0-9]+}/update", updateItem).Methods("POST")
	deleteItem := requireUserMw.ApplyFn(itemsController.Delete)
	r.HandleFunc("/items/{id:[0-9]+}/delete", deleteItem).Methods("POST")
	indexItem := requireUserMw.ApplyFn(itemsController.Index)
	r.Handle("/items", indexItem).Methods("GET")

	// COOKIE TEST
	r.HandleFunc("/cookietest", usersController.CookieTest).Methods("GET")

	// 404
	r.NotFoundHandler = staticController.Error404View

	// Server start
	http.ListenAndServe(":3000", r)
}
