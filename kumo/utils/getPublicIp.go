package utils

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// Returns the public IP of the machine running the program.
// Returns an empty string and an error if the call fails.
//
// Makes an external call.
//
// Example:
//
//	() -> ("123.456.789.012", nil)
func GetPublicIp() (ip string, err error) {
	var (
		response      *http.Response
		bytesResponse []byte
	)
	
	// Send GET request to retrieve public IP
	if response, err = http.Get("https://api.ipify.org?format=text"); err != nil {
		return "", errors.Wrap(err, "Error occurred while getting public IP")
	}
	defer func() {
		if errClosingBody := response.Body.Close(); err != nil {
			err = errors.Wrap(errClosingBody, "Error occurred while closing response body")
		}
	}()

	// Read the response body
	if bytesResponse, err = io.ReadAll(response.Body); err != nil {
		return "", errors.Wrap(err, "Error occurred while reading response body")
	}

	// Convert the response body to a string
	ip = string(bytesResponse)

	return
}
