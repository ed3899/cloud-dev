package utils

import (
	"os"

	"github.com/pkg/errors"
)

// Checks if a directory exists at the specified path. Preferably use absolute paths.
func DirExist(dirpath string) (bool, error) {
	// Get the file info of the directory
	_, err := os.Stat(dirpath)

	switch {
	// If the error indicates that the directory does not exist
	case os.IsNotExist(err):
		return false, nil // Return false to indicate that the directory does not exist
	// If the error indicates that the directory already exists
	case os.IsExist(err):
		return true, nil // Return true to indicate that the directory exists
	// If any other error occurred while checking the existence of the directory
	case err != nil:
		err = errors.Wrap(err, "Error occurred while checking existence of directory")
		return false, err // Return false and the wrapped error
	// If no error occurred during the check
	default:
		return true, nil // Return true to indicate that the directory exists
	}
}

// Checks if a directory does not exist at the specified path. Preferably use absolute paths.
func DirNotExist(dirpath string) (bool, error) {
	exists, err := DirExist(dirpath)
	if err != nil {
		return false, nil
	}
	return !exists, nil
}
