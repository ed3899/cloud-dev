package utils

import (
	"log"
	"os"

	"github.com/pkg/errors"
)

func GetCWD() (string, error) {
	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		err = errors.Wrap(err, "failed to get current directory")
		log.Printf("there was an error getting the current directory: %v", err)
		return "", err
	}

	// Return the current directory
	return dir, nil
}
