package utils

import (
	"log"
	"os"

	"github.com/pkg/errors"
)

// Returns the absolute path of the current working directory and any error that occurred.
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