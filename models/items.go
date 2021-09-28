package models

import (
	"github.com/jinzhu/gorm"
)

const (
	ErrUserIDRequired modelError = "models: user ID is required"
	ErrNameRequired	  modelError = "models: name is required"
)

type Item struct {
	gorm.Model
	UserID 		uint   `gorm:"not_null;index"`
	Name 		string `gorm:"not_null"`
	Price		string `gorm:"not_null"`
	Description string
}

type ItemService interface {
	ItemDB
}

type ItemDB interface {
	ByID(id uint) (*Item, error)
	Create(item *Item) error
}

type itemGorm struct {
	db *gorm.DB
}

type itemService struct {
	ItemDB
}

type itemValidator struct {
	ItemDB
}

var _ ItemDB = &itemGorm{}

type itemValFn func(*Item) error

func runItemValFns(item *Item, fns ...itemValFn) error {
	for _, fn := range fns {
		if err := fn(item); err != nil {
			return err
		}
	}

	return nil
}

func NewItemService(db *gorm.DB) ItemService {
	return &itemService{
		ItemDB: &itemValidator{
			ItemDB: &itemGorm{
				db: db,
			},
		},
	}
}

// Create will create the provided item and backfill data like the ID, CreatedAt, UpdatedAt fields
func (ig *itemGorm) Create(item *Item) error {
	return ig.db.Create(item).Error
}
func (iv *itemValidator) Create(item *Item) error {
	err := runItemValFns(item,
		iv.userIDRequired,
		iv.nameRequired)
	if err != nil {
		return err
	}

	return iv.ItemDB.Create(item)
}

func (ig *itemGorm) ByID(id uint) (*Item, error) {
	var item Item
	db := ig.db.Where("id = ?", id)
	err := first(db, &item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (iv *itemValidator) userIDRequired(i *Item) error {
	if i.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (iv *itemValidator) nameRequired(i *Item) error {
	if i.Name == "" {
		return ErrNameRequired
	}
	return nil
}
