package models

import (
	"github.com/jinzhu/gorm"
)

type Item struct {
	gorm.Model
	Name 		string `gorm:"not_null"`
	Price		string `gorm:"not_null"`
	Description string
}

type ItemService interface {
	ItemDB
}

type ItemDB interface {
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
