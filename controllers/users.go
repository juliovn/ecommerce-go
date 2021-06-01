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
		us:		 us,
	}
}

type Users struct {
	NewView *views.View
	us		*models.UserService
}

type SignupForm struct {
	Name		string `schema:"name"`
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
