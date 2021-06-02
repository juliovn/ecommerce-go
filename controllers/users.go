package controllers

import (
	"ecommerce/models"
	"ecommerce/views"
	"fmt"
	"net/http"
)

func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("base", "users/new"),
		LoginView: views.NewView("base", "users/login"),
		us:		 us,
	}
}

type Users struct {
	NewView 	*views.View
	LoginView	*views.View
	us			*models.UserService
}

type SignupForm struct {
	Name		string `schema:"name"`
	Email		string `schema:"email"`
	Password	string `schema:"password"`
}

type LoginForm struct {
	Email		string `schema:"email"`
	Password	string `schema:"password"`
}

func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Create is used to process the signup form when user creates a new account
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user := models.User{
		Name: form.Name,
		Email: form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User is", user)
}

// Login is used to process the login form and will attempt to log in a user
//
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user, err := u.us.Authenticate(form.Email, form.Password)

	switch err {
	case models.ErrNotFound:
		fmt.Fprintln(w, "Invalid username")
	case models.ErrInvalidPassword:
		fmt.Fprintln(w, "Invalid password")
	case nil:
		fmt.Fprintln(w, user)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
