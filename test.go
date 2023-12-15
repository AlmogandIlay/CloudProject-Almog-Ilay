package main

import (
	"fmt"
	"server/authentication"
)

func main() {
	fmt.Println("Starting Testing")
	ml, err := authentication.InitializeLoginManager()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("\nValid Signup:")
	errs := ml.Signup("almog", "87654321", "almog1@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Println("Logged Users after a valid signup:", ml.GetLoggedUsers())
	fmt.Println("\nInvalid Signup: Short Username")
	errs = ml.Signup("A", "87654321", "almog2@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Println("Logged users after an invalid signup: ", ml.GetLoggedUsers())
}
