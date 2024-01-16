package Filesystem

func (user *LoggedUser) ChangeDir(newPath string) error {
	return user.setPath(newPath)
}

// root/filesys/fol/ilay.txt
