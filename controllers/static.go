package controllers

import "ecommerce/views"

func NewStatic() *Static {
	return &Static{
		HomeView: views.NewView("base", "views/static/home.gohtml"),
		ContactView: views.NewView("base", "views/static/contact.gohtml"),
		Error404View: views.NewView("base", "vies/static/error404.gohtml"),
	}
}


type Static struct {
	HomeView 		*views.View
	ContactView 	*views.View
	Error404View	*views.View
}