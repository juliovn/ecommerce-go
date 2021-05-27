package models

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrNotFound is returned when a resouece cannot be found in the database
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is provided to a method like Delete
	ErrInvalidID = errors.New("models: ID provided was invalid")
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
// If there is another error, return an error with more information
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ByEmail looks up a user with the given email address and returns that user
// If the user is found, return a nil error
// If the user is not found, return ErrNotFound
// If there is another error, return error with more information
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// Update will update the provided user with all of the data in the provided user object
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

// Create will create the provided user and backfill data like the ID, CreatedAt, UpdatedAt fields
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// Delete will delete the user with the provided ID
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

// DestructiveReset drops the user table and rebuilds it
func (us *UserService) DestructiveReset() error {
	err := us.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return err
	}
	return us.AutoMigrate()
}

// AutoMigrate will attempt to automatically migrate the users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}


// Closes the UserService database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// first will query using the provided gorm.DB and it will get the first
// item returned and place it into dst. if nothing is found in the query
// it will return ErrNotFound
func first (db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}
