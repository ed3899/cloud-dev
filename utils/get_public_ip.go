package utils

import (
	"io"
	"log"
	"net/http"

	"github.com/samber/oops"
)

type GetPublicIpF func() (ip string, err error)

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
		oopsBuilder = oops.Code("get_public_ip_failed")
		URL         = "https://api.ipify.org?format=text"

		response      *http.Response
		bytesResponse []byte
	)

	// Send GET request to retrieve public IP
	if response, err = http.Get(URL); err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while sending GET request to '%s'", URL)
		return
	}
	defer func(response *http.Response) {
		if err := response.Body.Close(); err != nil {
			log.Fatalf(
				"%+v",
				oopsBuilder.
					With("responseStatusCode", response.StatusCode).
					Wrapf(err, "Error occurred while closing response body"),
			)
		}
	}(response)

	// Read the response body
	if bytesResponse, err = io.ReadAll(response.Body); err != nil {
		err = oopsBuilder.
			With("responseStatusCode", response.StatusCode).
			Wrapf(err, "Error occurred while reading response body")
		return
	}

	// Convert the response body to a string
	ip = string(bytesResponse)

	return
}
