package FileSystem

type Folder struct {
	Name    string    // File name without extension
	Path    string    // File path
	Files   []File    // Files contain in the folder
	Folders []*Folder // Slice of pointers to Folder structs
}
