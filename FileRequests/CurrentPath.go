package FileRequestsManager

import "fmt"

var (
	CurrentPath string
)

func InitializeCurrentPath() {
	CurrentPath = "Root:\\"
}

func PrintCurrentPath() {
	fmt.Print(CurrentPath)
}

func IsCurrentPathInitialized() bool {
	return CurrentPath != ""
}

func SetCurrentPath(path string) {
	CurrentPath = path
}
