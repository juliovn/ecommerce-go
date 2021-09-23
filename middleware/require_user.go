package middleware

import (
	"ecommerce/models"
	"fmt"
	"net/http"
)

type RequireUser struct {
	models.UserService
}

// ApplyFn will return an http.HandlerFunc that will
// check to see if a user is logged in and then either call
// next(w, r) if they are, or redirect them to login page if not
func (mw *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	// return a dynamically created func(http.ResponseWriter, *http.Request)
	// but also need to convert it into a http.HandlerFunc
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
		}

		fmt.Println("User found:", user)
		next(w, r)
	})
}

func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.ApplyFn(next.ServeHTTP)
}