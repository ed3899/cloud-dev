package utils

import (
	"log"
	"os"
)

func DependenciesToBeDownloaded(dp []*Dependency) []*Dependency {
	for _, d := range dp {
		if FileExists(d.ExtractionPath) {
			d.Present = true
		} else {
			d.Present = false
		}
	}
	return dp
}

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if err == nil {
		return true // File exists
	}
	if os.IsNotExist(err) {
		log.Printf("File '%s' does not exist", filepath)
		return false // File does not exist
	}
	return false // Error occurred while checking file existence
}
