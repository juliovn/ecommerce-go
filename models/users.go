package models

import (
	"ecommerce/hash"
	"ecommerce/rand"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

const (
	hmacSecretKey = "6hMXiDnQoH8Ec7KQ"
	userPwPepper  = "6J53aQog5tQPaJPT"

	passwordMinLength = 4
)

var (
	// ErrNotFound is returned when a resouece cannot be found in the database
	ErrNotFound modelError = "models: resource not found"

	// ErrIdInvalid is returned when an invalid ID is provided to a method like Delete
	ErrIdInvalid modelError = "models: ID provided was invalid"

	// ErrPasswordIncorrect is returned when an invalid password is used when attempting to authenticate a user
	ErrPasswordIncorrect modelError = "models: incorrect password provided"

	// ErrEmailRequired is returned when a username is not provided when creating a user
	ErrEmailRequired modelError = "models: username is required"

	// ErrEmailTaken is returned when a username already exists on database
	ErrEmailTaken modelError = "models: email address is already taken"

	// ErrPasswordTooShort is returned when a user tries to set a password that is lower
	// than passwordMinLength const
	ErrPasswordTooShort modelError = "models: password is too short"

	// ErrPasswordRequired is returned when a Create is attempted without a user password
	ErrPasswordRequired modelError = "models: password is required"
)

// Public and Error wrap errors for display on the frontend
type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

type UserDB interface {
	// Methods for querying single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

type userGorm struct {
	db *gorm.DB
}

var _ UserDB = &userGorm{}
var _ UserService = &userService{}

type userService struct {
	UserDB
}

type userValidator struct {
	UserDB
	hmac hash.HMAC
}

type userValFn func(*User) error

type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

func NewUserService(db *gorm.DB) UserService {
	ug := &userGorm{db}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)
	return &userService{
		UserDB: uv,
	}
}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB: udb,
		hmac:   hmac,
	}
}

// ByID will look up a user with the provided ID
// If the user is found, return nil error
// If the user is not found, return ErrNotFound
// If there is another error, return an error with more information
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
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
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	err := runUserValFns(&user, uv.normalizeEmail)
	if err != nil {
		return nil, err
	}

	return uv.UserDB.ByEmail(user.Email)
}

// ByRemember looks up a user with the given remember token and returns that user.
// This method will handle hashing the token
// Errors are the same as ByEmail
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}

	if err := runUserValFns(&user, uv.hmacRemember); err != nil {
		return nil, err
	}

	return uv.UserDB.ByRemember(user.RememberHash)
}

// Update will update the provided user with all of the data in the provided user object
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}
func (uv *userValidator) Update(user *User) error {
	err := runUserValFns(user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.hmacRemember,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailIsAvailable)
	if err != nil {
		return err
	}

	return uv.UserDB.Update(user)
}

// Create will create the provided user and backfill data like the ID, CreatedAt, UpdatedAt fields
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}
func (uv *userValidator) Create(user *User) error {
	err := runUserValFns(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setRememberIfUnset,
		uv.hmacRemember,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailIsAvailable)
	if err != nil {
		return err
	}

	return uv.UserDB.Create(user)
}

// Delete will delete the user with the provided ID
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}
func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFns(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}

	return uv.UserDB.Delete(id)
}

// Authenticate can be used to authenticate a user with the provided username and password
// If the username provided is invalid, this will return ErrNotFound
// If the password is invalid, this will return ErrPasswordIncorrect
// If the email and password are both valid, this will return user, nil
// Otherwise if another error is encountered this will return nil, error
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash),
		[]byte(password+userPwPepper))

	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrPasswordIncorrect
	default:
		return nil, err
	}
}

// first will query using the provided gorm.DB and it will get the first
// item returned and place it into dst. if nothing is found in the query
// it will return ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// bcryptPassword will hash a user's password with an app-wide
// pepper and bcrypt, which salts
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

// runUserValFns will take a pointer to a user and a variadic number of validation functions
// and run each one sequentially, returning an error where it is encountered
func runUserValFns(user *User, functions ...userValFn) error {
	for _, fn := range functions {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

// hmacRemember is a validation function that will check to see if remember token is present
// and if so will hash using HMAC
func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

// setRememberIfUnset will prevent a user being saved on the database without a remember token
func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

// idGreaterThan will ensure a valid ID for delete
func (uv *userValidator) idGreaterThan(n uint) userValFn {
	return userValFn(func(user *User) error {
		if user.ID <= n {
			return ErrIdInvalid
		}
		return nil
	})
}

// normalizeEmail performs trimming and transform email to lowercase
func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

// requireEmail makes sure an email has been provided
func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}

	return nil
}

// emailIsAvailable will query the database for a user with provided email
// if user is not found, return nil continuing the chain
// if query fails for some reason, return error
// if the user.ID provided is the same as the user.ID of the queried object, that means this is taken so return ErrEmailTaken
func (uv *userValidator) emailIsAvailable(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		return nil
	}

	if err != nil {
		return err
	}

	if user.ID != existing.ID {
		return ErrEmailTaken
	}

	return nil
}

// passwordMinLength will make sure password meets minimum length
func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 4 {
		return ErrPasswordTooShort
	}

	return nil
}

// passwordRequired makes sure password field is not empty
func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

// passwordHashRequired is the same as passwordRequired but for the hash
func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}
	return nil
}
