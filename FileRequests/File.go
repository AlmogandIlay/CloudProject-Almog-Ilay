package FileRequestsManager

type file struct {
	Name string `json:"name"` // File name (including its extension)
	Path string `json:"path"` // File's path in the Cloud
	Size uint32 `json:"size"` // File's size in bytes
}

// Creates a new file struct with the given parameters
func newFile(name string, path string, size uint32) file {
	return file{
		Name: name,
		Path: path,
		Size: size,
	}
}
