package models

import "github.com/jinzhu/gorm"

type Item struct {
	gorm.Model
	Name string `gorm:"not_null"`
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

// Create will create the provided item and backfill data like the ID, CreatedAt, UpdatedAt fields
func (ig *itemGorm) Create(item *Item) error {
	// TODO implement this function
	return nil
}
