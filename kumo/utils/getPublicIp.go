package utils

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// Returns the public IP of the machine running the program. Makes an external call.
func GetPublicIp() (ip string, err error) {
	// Send GET request to retrieve public IP
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting public IP")
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	bytesResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while reading response body")
		return "", err
	}

	// Convert the response body to a string
	ip = string(bytesResp)

	 // Return the public IP and no error
	return ip, nil
}
