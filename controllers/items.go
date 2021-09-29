package controllers

import (
	"ecommerce/context"
	"ecommerce/models"
	"ecommerce/views"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	ShowItem = "show_item"
)

type Items struct {
	New 		*views.View
	ShowView	*views.View
	EditView	*views.View
	IndexView	*views.View
	is			models.ItemService
	r			*mux.Router
}

type ItemForm struct {
	Name		string `schema:"name"`
	Description	string `schema:"description"`
	Price		string `schema:"price"`
}

func NewItems(is models.ItemService, r *mux.Router) *Items {
	return &Items{
		New:    	views.NewView("base", "items/new"),
		ShowView: 	views.NewView("base", "items/show"),
		EditView:	views.NewView("base", "items/edit"),
		IndexView: 	views.NewView("base", "items/index"),
		is: 		is,
		r: 			r,
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


	url, err := i.r.Get(ShowItem).URL("id", strconv.Itoa(int(item.ID)))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	http.Redirect(w, r, url.Path, http.StatusFound)
}

// GET /items/:id
func (i *Items) Show(w http.ResponseWriter, r *http.Request) {
	item, err := i.itemByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if item.UserID != user.ID {
		http.Error(w, "You do not have permissions to view this item", http.StatusForbidden)
		return
	}

	var vd views.Data
	vd.Yield = item
	i.ShowView.Render(w, vd)
}

// GET /items/:id/edit
func(i *Items) Edit(w http.ResponseWriter, r *http.Request) {
	item, err := i.itemByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if item.UserID != user.ID {
		http.Error(w, "You do not have permissions to edit this item", http.StatusForbidden)
		return
	}

	var vd views.Data
	vd.Yield = item
	i.EditView.Render(w, vd)
}

func (i *Items) itemByID(w http.ResponseWriter, r *http.Request) (*models.Item, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusNotFound)
	}

	item, err := i.is.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Item not found", http.StatusNotFound)
		default:
			http.Error(w, "Something went wrong on item ID lookup", http.StatusInternalServerError)
		}
		return nil, err
	}

	return item, nil
}

// POST /items/:id/update
func (i *Items) Update(w http.ResponseWriter, r *http.Request) {
	item, err := i.itemByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if item.UserID != user.ID {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	var vd views.Data
	vd.Yield = item
	var form ItemForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		i.EditView.Render(w, vd)
		return
	}

	item.Name = form.Name
	item.Price = form.Price
	item.Description = form.Description
	err = i.is.Update(item)
	if err != nil {
		vd.SetAlert(err)
	} else {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Item successfully updated",
		}
	}

	i.EditView.Render(w, vd)
}

// POST /items/:id/delete
func (i *Items) Delete(w http.ResponseWriter, r *http.Request) {
	item, err := i.itemByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if item.UserID != user.ID {
		http.Error(w, "You don't have permissions to delete this item", http.StatusForbidden)
		return
	}

	var vd views.Data
	err = i.is.Delete(item.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = item
		i.EditView.Render(w, vd)
		return
	}

	fmt.Fprintln(w, "successfully deleted")
}

// GET /items
func (i *Items) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	items, err := i.is.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong with ByUserID lookup", http.StatusInternalServerError)
		return
	}
	var vd views.Data
	vd.Yield = items
	i.IndexView.Render(w, vd)
}
