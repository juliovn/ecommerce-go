package controllers

import "ecommerce/views"

func NewStatic() *Static {
	return &Static{
		HomeView:     views.NewView("base", "static/home"),
		Error404View: views.NewView("base", "static/error404"),
	}
}

type Static struct {
	HomeView     *views.View
	Error404View *views.View
}
