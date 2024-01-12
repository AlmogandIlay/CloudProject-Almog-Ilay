package Filesystem

type Folder struct {
	Name    string // name without extension
	Path    string
	Files   []File
	Folders []*Folder // Slice of pointers to Folder structs
}
