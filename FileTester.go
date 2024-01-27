package main

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

	id := uint32(871)

	user, err := FileSystem.NewLoggedUser(id)
	if err != nil {
		fmt.Println("test: " + err.Error())
		os.Exit(1)
	}

	printUser(user)
	err = user.CreateFile("C:\\CloudDrive\\871\\f1\\file.txt")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	printUser(user)

}
