package controllers

import (
	"ecommerce/context"
	"ecommerce/models"
	"ecommerce/views"
	"fmt"
	"net/http"
)

type Items struct {
	New 	*views.View
	is		models.ItemService
}

type ItemForm struct {
	Name		string `schema:"name"`
	Description	string `schema:"description"`
	Price		string `schema:"price"`
}

func NewItems(is models.ItemService) *Items {
	return &Items{
		New:    views.NewView("base", "items/new"),
		is: 	is,
	}
}

// POST /items
func (i *Items) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form ItemForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		i.New.Render(w, vd)
		return
	}

	user := context.User(r.Context())

	item := models.Item{
		Name: form.Name,
		Description: form.Description,
		Price: form.Price,
		UserID: user.ID,
	}

	if err := i.is.Create(&item); err != nil {
		vd.SetAlert(err)
		i.New.Render(w, vd)
		return
	}
	fmt.Fprintln(w, item)
}
