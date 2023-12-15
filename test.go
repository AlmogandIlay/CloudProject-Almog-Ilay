package main

import (
	"fmt"
	"server/authentication"
)

func main() {

	user, errs := authentication.NewUser("AA", "8764321", "almog2@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Print(user)
	/*ml, err := authentication.InitializeLoginManager()

	if err != nil {
		fmt.Print(err.Error())
	}
	_ = ml.Signup("ilmog", "12345678", "ilmog5@gmail.com")
	fmt.Println("Valid Signup:")
	errs := ml.Signup("almog", "87654321", "almog1@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Println("Logged Users after a valid signup:", ml.GetLoggedUsers())
	fmt.Println("\nInvalid Signup: Short Username")
	errs := ml.Signup("AA", "8764321", "almog2@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Println("Logged users after an invalid signup: ", ml.GetLoggedUsers())*/

}
