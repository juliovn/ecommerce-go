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
	ByUserID(userID uint) ([]Item, error)
	Create(item *Item) error
	Update(item *Item) error
	Delete(id uint) error
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

func (ig *itemGorm) Update(item *Item) error {
	return ig.db.Save(item).Error
}
func (iv *itemValidator) Update(item *Item) error {
	err := runItemValFns(item,
		iv.userIDRequired,
		iv.nameRequired)
	if err != nil {
		return err
	}
	return iv.ItemDB.Update(item)
}

func (ig *itemGorm) Delete(id uint) error {
	item := Item{ Model: gorm.Model{ ID:id } }
	return ig.db.Delete(&item).Error
}
func (iv *itemValidator) Delete(id uint) error {
	var item Item
	item.ID = id
	if err := runItemValFns(&item, iv.nonZeroID); err != nil {
		return err
	}
	return iv.ItemDB.Delete(item.ID)
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

func (ig *itemGorm) ByUserID(userID uint) ([]Item, error) {
	var items []Item
	db := ig.db.Where("user_id = ?", userID)
	if err := db.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
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

func (iv *itemValidator) nonZeroID(item *Item) error {
	if item.ID <= 0 {
		return ErrIdInvalid
	}
	return nil
}