package Helper

import (
	"os"
	"path/filepath"
	"strings"
)

// Returns whether the given path exists. Returns error if the check gone wrong
func IsPathExists(path string) (bool, error) {
	if path != "" { // If path has any value, otherwise returns true automatically
		_, err := os.Stat(path)
		if err == nil { // If path exists
			return true, nil
		}
		if os.IsNotExist(err) { // If path not exists
			return false, nil
		}
		// If the check gone wrong
		return false, err
	}
	return true, nil
}

// Counting contents (files and folders) in a given folderpath
func CountContents(folderPath string) (uint, uint, error) {
	folderPath = strings.ReplaceAll(folderPath, "'", "") // Removes ' ' between the path if exists
	var filesCount, folderCount uint

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		// for each content in the directory
		if err != nil {
			return err
		}
		if info.IsDir() { // If it's a directory
			folderCount++
		} else { // If it's a file
			filesCount++
		}

		return nil
	})
	// Substract 1 from foldercount to remove the counting of the folderPath itself
	return filesCount, folderCount - 1, err
}

// Find if the parameters are closed with ' in the amount of the given closedCount (for example, given TwoClosedPath: 'filename1' 'filename2')
func IsQuoted(command_arguments []string, closedCount AmountOfPaths) bool {
	arguments := strings.Join(command_arguments, " ")
	var counter AmountOfPaths = 0
	for _, char := range arguments {
		if char == '\'' {
			counter++
		}
	}
	return counter == closedCount
}

// Find the path of a file if it's closed with ' with the given amount of closed paths and specific filename index.
// Returns the string path with enclouse '
func FindPath(command_arguments []string, index int, closedCount AmountOfPaths) string {
	arguments := strings.Join(command_arguments, " ") // Convert to string
	var name string
	var counter int // counting the amount of '
	for _, char := range arguments {
		if char == '\'' { // If close ' char found
			counter++
		}
		name += string(char)    // append char to string
		if counter >= index+2 { // If 2/4 ' close chars have been found
			break // Stop searching
		}
	}
	if index == SecondNameParameter { // If the second path has been chosen
		parts := strings.Split(name, "'")
		// If the client provided 2 paths (two single quotes), save the second part
		if len(parts) >= int(TwoCloudPaths) {
			name = parts[secondPathIndex]
		} else {
			name = ""
		}
		//name = "'" + strings.Join((strings.Split(name, "'")[1])) + "'" // save the second path only. (int[closedCount]-1) to access the last index of the closedCount path
	}
	return name
}

// For cases when the first path is quoted but the second doesn't, return the second Non-quoted second path
func ReturnNonQuotedSecondPath(command_arguments []string) string {
	arguments := strings.Join(command_arguments, " ") // Convert the comman arguments to string
	lastIndex := strings.LastIndex(arguments, "'")    // Find the last time the ' index encloused has been used in the string and append 1 to avoid it completely
	if lastIndex+1 <= len(arguments) {                // If Second path hasn't been specified
		return ""
	}
	return arguments[lastIndex+1:]
}
