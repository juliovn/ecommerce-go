package controllers

import (
	"ecommerce/models"
	"ecommerce/views"
)

type Items struct {
	New 	*views.View
	is		models.ItemService
}

func NewItems(is models.ItemService) *Items {
	return &Items{
		New:    views.NewView("base", "items/new"),
		is: 	is,
	}
}

