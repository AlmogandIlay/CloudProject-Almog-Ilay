package main

import (
	"fmt"
	"server/authentication"
)

func main() {
	fmt.Println("Starting Testing\n")
	/*var username string
	var password string
	var email string
	fmt.Println("Enter your username: ")
	fmt.Scanln(&username)
	fmt.Println("Enter your password: ")
	fmt.Scanln(&password)
	fmt.Println("Enter your email: ")
	fmt.Scanln(&email)
	*/
	ml, err := authentication.InitializeLoginManager()

	if err != nil {
		fmt.Print(err.Error())
	}
	//_ = ml.Signup("ilmog", "12345678", "ilmog5@gmail.com")
	fmt.Println("Valid Signup:")
	errs := ml.Signup("almog", "87654321", "almog1@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Println("Logged Users after a valid signup:", ml.GetLoggedUsers())
	fmt.Println("\nInvalid Signup: Short Username")
	errs = ml.Signup("A", "8764321", "almog2@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Println("Logged users after an invalid signup: ", ml.GetLoggedUsers())

}
