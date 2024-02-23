package FileRequestsManager

type file struct {
	Name string // File name (including its extension)
	Path string // File's path in the Cloud
	Size uint32 // File's size in bytes
}

func newFile(name string, path string, size uint32) file {
	return file{
		Name: name,
		Path: path,
		Size: size,
	}
}
