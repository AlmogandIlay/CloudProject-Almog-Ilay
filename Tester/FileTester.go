package Tester

import (
	"CloudDrive/FileSystem"
	"fmt"
	"os"
	//"path/filepath"
	//"strconv"
)

func printUser(user *FileSystem.LoggedUser) {
	fmt.Printf("test: user{id: %d, path: %s}\n", user.UserID, user.CurrentPath)
}
func main() {

	id := uint32(872)

	user, err := FileSystem.NewLoggedUser(id)
	if err != nil {
		fmt.Println("test: " + err.Error())
		os.Exit(1)
	}

	printUser(user)

	path, err := user.ChangeDirectory("C:\\CloudDrive\\872\\ilay\\aa")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(4)
	}

	fmt.Println(path)

}
