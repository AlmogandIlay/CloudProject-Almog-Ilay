package main

import (
	"fmt"
	"server/authentication"
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
	manager, err := authentication.NewLoginManager()
	if err != nil {
		fmt.Println(err.Error())
	}
	manager.Signup(username, password, email)
	fmt.Println(manager.GetLoggedUsers())

}
