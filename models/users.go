package models

import (
	"ecommerce/hash"
	"ecommerce/rand"
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

const hmacSecretKey = "6hMXiDnQoH8Ec7KQ"

var (
	// password pepper string
	userPwPepper = "6J53aQog5tQPaJPT"

	// ErrNotFound is returned when a resouece cannot be found in the database
	ErrNotFound = errors.New("models: resource not found")

	// ErrInvalidID is returned when an invalid ID is provided to a method like Delete
	ErrInvalidID = errors.New("models: ID provided was invalid")

	// ErrInvalidPassword is returned when an invalid password is used when attempting to authenticate a user
	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

type UserService struct {
	db 		*gorm.DB
	hmac	hash.HMAC
}

type User struct {
	gorm.Model
	Name			string
	Email			string `gorm:"not null;unique_index"`
	Password		string `gorm:"-"`
	PasswordHash	string `gorm:"not null"`
	Remember		string `gorm:"-"`
	RememberHash	string `gorm:"not null;unique_index"`
}

func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	db.LogMode(true)

	hmac := hash.NewHMAC(hmacSecretKey)

	return &UserService{
		db: db,
		hmac: hmac,
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

// ByRemember looks up a user with the given remember token and returns that user.
// This method will handle hashing the token
// Errors are the same as ByEmail
func (us *UserService) ByRemember(token string) (*User, error) {
	var user User
	rememberHash := us.hmac.Hash(token)
	err := first(us.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update will update the provided user with all of the data in the provided user object
func (us *UserService) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = us.hmac.Hash(user.Remember)
	}

	return us.db.Save(user).Error
}

// Create will create the provided user and backfill data like the ID, CreatedAt, UpdatedAt fields
func (us *UserService) Create(user *User) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password + userPwPepper), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}

	user.RememberHash = us.hmac.Hash(user.Remember)

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

// Authenticate can be used to authenticate a user with the provided username and password
// If the username provided is invalid, this will return ErrNotFound
// If the password is invalid, this will return ErrInvalidPassword
// If the email and password are both valid, this will return user, nil
// Otherwise if another error is encountered this will return nil, error
func (us *UserService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash),
		[]byte(password + userPwPepper))

	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrInvalidPassword
	default:
		return nil, err
	}
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
