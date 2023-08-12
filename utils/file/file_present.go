package file

import (
	"io/fs"
)

// Returns true if the file at the given os.File exists, false otherwise.
func FilePresentWith(
	osStat func(string) (fs.FileInfo, error),
	osIsNotExist func(error) bool,
) ForPath {
	var (
		err error
	)

	forPath := func(absolutePath string) bool {
		// Check if the file at the given absolute path exists
		_, err = osStat(absolutePath)
		// If there is no error, it means the file exists
		if err == nil {
			return true
		}

		// If the error indicates that the file does not exist, return false
		if osIsNotExist(err) {
			return false
		}

		// If there is any other type of error, assume the file does not exist
		return false
	}

	return forPath
}

type ForPath func(string) bool
