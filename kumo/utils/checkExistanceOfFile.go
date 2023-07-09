package utils

import (
	"os"
)

func FilePresent(absolutePath string) bool {
	_, err := os.Stat(absolutePath)
	if err == nil {
		return true // File exists
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func FileNotPresent(absolutePath string) bool {
	return !FilePresent(absolutePath)
}
