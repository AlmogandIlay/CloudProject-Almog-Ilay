package main

import (
	"fmt"
	"server/authentication"
)

func main() {
	fmt.Println("Hey")
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
	ml, err := authentication.NewLoginManager()

	if err != nil {
		fmt.Print(err.Error())
	}
	_ = ml.Signup("ilmog", "12345678", "ilmog5@gmail.com")
	fmt.Print(ml.GetLoggedUsers())
}
