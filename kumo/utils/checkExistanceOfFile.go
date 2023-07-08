package utils

import (
	"os"
)

func DependencyPresent(absolutePath string) bool {
	_, err := os.Stat(absolutePath)
	if err == nil {
		return true // File exists
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func DependencyNotPresent(absolutePath string) bool {
	return !DependencyPresent(absolutePath)
}
