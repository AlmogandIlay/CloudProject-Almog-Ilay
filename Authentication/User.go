package Authentication

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func Signup(username string, password string, email string) User {
	return User{Username: username, Password: password, Email: email}
}

func Signin(username string, password string) User {
	return User{Username: username, Password: password, Email: ""}
}
