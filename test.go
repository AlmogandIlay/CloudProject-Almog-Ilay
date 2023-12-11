package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hey")
	var username string
	var password string
	var email string
	fmt.Println("Enter your username: ")
	fmt.Scanln(&username)
	fmt.Println("Enter your password: ")
	fmt.Scanln(&password)
	fmt.Println("Enter your email: ")
	fmt.Scanln(&email)
	user, err := NewUser(username, password, email)
	if err != nil {
		panic(err)
	}
	fmt.Println((user.username))
}
