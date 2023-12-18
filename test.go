package main

import (
	"fmt"
	"server/authentication"
)

func main() {
	fmt.Println("-----Starting Testing-----")
	fmt.Println("-----Starting Signup Tests-----")
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
	errs = ml.Signup("JFKFISOMADFJKSDFJKSDFKJSKJSDFKLJSDFKLJSDKJDSF", "87654321", "almog3@gmail.com")
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
	}
	fmt.Println("Logged users after an exact Email Address Signup", ml.GetLoggedUsers())

	fmt.Println("\n---------Multiple Errors:----------")

	fmt.Println("\nInvalid Signup: Short Username + Short Password")
	errs = ml.Signup("A", "B", "almog@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Println("Logged users after a short username + short password Signup", ml.GetLoggedUsers())

	fmt.Println("\nInvalid Signup: Short Username + Long Password")
	errs = ml.Signup("A", "dfsjhgsdjfsdfjhksdfkhsdkjhfdskhjfsdkjfhdskfjhsdjfhf", "almog@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Println("Logged users after a short Username + long Password Signup", ml.GetLoggedUsers())

	fmt.Println("\nInvalid Signup: Long Username + Short Password")
	errs = ml.Signup("RedJet999", "A", "almog@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Println("Logged users after a short username + long password Signup", ml.GetLoggedUsers())

	fmt.Println("\nInvalid Signup: Empty Username + Empty Password + Empty mail")
	errs = ml.Signup("", "", "")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Println("Logged users after an empty username, password, email Signup", ml.GetLoggedUsers())

	fmt.Println("\nInvalid Signup: Long Username + Long Password")
	errs = ml.Signup("sdfjhdsfjhsdfjjhdsjdsf", "sdfjkdsjkdsfjksdfjksadjfkasdj", "almog@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	fmt.Println("Logged users after a long username + long password Signup", ml.GetLoggedUsers())

	fmt.Println("\n-------------Testing Login Tests-------------")

	fmt.Println("\nValid Login:")
	err = ml.Login("almog", "87654321")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("\nInvalid Login: Incorrect username")
	err = ml.Login("correct", "87654321")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("\nInvalid Login: Incorrect Password")
	err = ml.Login("almog", "correct")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("\nInvalid Login: long username")
	err = ml.Login("almogggggggggggggggggggggggggggggg", "87654321")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("\nInvalid Login: long password")
	err = ml.Login("almog", "87654321111111111111111111111111111111111111111111")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("\nInvalid Login: Incorrect username and password")
	err = ml.Login("correct", "incorrect")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("\n---------Testing Logout-----------")
	fmt.Println("Creating a few users...")

	errs = ml.Signup("Alex", "Almog123", "almog6@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	errs = ml.Signup("Loan", "Almog123", "almog7@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	errs = ml.Signup("Corin", "Almog123", "almog8@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	errs = ml.Signup("Spark", "Almog123", "almog9@gmail.com")
	for _, err := range errs {
		fmt.Println(err.Error())
	}

	fmt.Println("Current Registered Users:\n", ml.GetLoggedUsers())

	fmt.Println("\n\nValid Logout")
	err = ml.Logout("Loan")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Registered Users:", ml.GetLoggedUsers())

	fmt.Println("\nInvalid Logout")
	err = ml.Logout("AmIExist?")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Registered Users:", ml.GetLoggedUsers())

}
