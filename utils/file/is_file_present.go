package file

import (
	"os"
)

// Returns true if the file at the given os.File exists, false otherwise.
func IsFilePresent(
	absolutePath string,
) bool {
	// Check if the file at the given absolute path exists
	_, err := os.Stat(absolutePath)
	// If there is no error, it means the file exists
	if err == nil {
		return true
	}

	// If the error indicates that the file does not exist, return false
	if os.IsNotExist(err) {
		return false
	}

	// If there is any other type of error, assume the file does not exist
	return false
}
