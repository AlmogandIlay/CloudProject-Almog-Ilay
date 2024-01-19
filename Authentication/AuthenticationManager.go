package Authentication

const (
	username_index = 0
	password_index = 1
	email_index    = 2
)

func HandleSignup(command_arguments []string) {
	user := Signup(command_arguments[username_index], command_arguments[password_index], command_arguments[email_index])
}
