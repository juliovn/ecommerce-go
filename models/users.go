package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a resouece cannot be found in the database
	ErrNotFound = errors.New("models: resource not found")
)

type UserService struct {
	db *gorm.DB
}

type User struct {
	gorm.Model
	Name	string
	Email	string `gorm:"not null;unique_index"`
}

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)

	return &UserService{
		db: db,
	}, nil
}

// ByID will look up a user with the provided ID
// If the user is found, return nil error
// If the user is not found, return ErrNotFound
// If there is another error, return an error with more information about what went wrong
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	err := us.db.Where("id = ?", id).First(&user).Error

	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
			return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Create will create the provided user and backfill data like the ID, CreatedAt, UpdatedAt fields


// DestructiveReset drops the user table and rebuilds it
func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

// Closes the UserService database connection
func (us *UserService) Close() error {
	return us.db.Close()
}
