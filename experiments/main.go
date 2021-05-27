package main

import (
	"fmt"

	"ecommerce/models"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "ecommerce_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.DestructiveReset()

	// Create a user
	user := models.User{
		Name: "Julius Caesar",
		Email: "julius.caesar@spqr.com",
	}

	if err := us.Create(&user); err != nil {
		panic(err)
	}

	// This will error because you DO NOT have a user with
	// this ID, but we will create one soon.
	foundUser, err := us.ByEmail("julius.caesar@spqr.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(foundUser)
}