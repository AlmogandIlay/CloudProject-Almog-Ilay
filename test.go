package main

import (
	"fmt"
	"server/authentication"
)

func main() {
	fmt.Println("-----Starting Testing-----")
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
	fmt.Println("Logged users after a short username signup: ", ml.GetLoggedUsers())

	fmt.Println("\nInvalid Signup: Long Username")
	errs = ml.Signup("JFKFISOMA", "87654321", "almog3@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Println("Logged users after a long username signup: ", ml.GetLoggedUsers())

	fmt.Println("\nInvalid Signup: Short Password")
	errs = ml.Signup("Dokter", "ksks", "almog4@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Println("Logged users after a short password signup:", ml.GetLoggedUsers())

	fmt.Println("\nInvalid Signup: Long Password")
	errs = ml.Signup("Jasper", "sdfjhdsfhfjhsdfjsdsd", "almog5@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Println("Logged users after a long password signup", ml.GetLoggedUsers())

	fmt.Println("\nInvalid Signup: Invalid email: @gmail.com")
	errs = ml.Signup("Rens", "Password123", "@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Println("\nInvalid Signup after a @gmail.com mail signup:", ml.GetLoggedUsers())

	fmt.Println("\nInvalid Signup: Invalid email: ...@")
	errs = ml.Signup("Ektiv", "Password456", "averyvalidmailaddress@")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Println("Logged users after a ...@ mail Signup", ml.GetLoggedUsers())

	fmt.Println("\nInvalid Signup: Same username")
	errs = ml.Signup("almog", "abeautifulpassword", "almog6@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Println("Logged users after an exact username Signup", ml.GetLoggedUsers())

	fmt.Println("\nInvalid Signup: Same Email Address")
	errs = ml.Signup("Candy", "Whatabeautifu", "almog1@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
		fmt.Print(err)
	}
	fmt.Println("Logged users after an exact Email Address Signup", ml.GetLoggedUsers())
	//fmt.Println("\nInvalid Sign")

}
