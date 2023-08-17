package ip

import (
	"io"
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
func GetPublicIp() (string, error) {
	oopsBuilder := oops.Code("get_public_ip_failed")
	URL := "https://api.ipify.org?format=text"

	// Send GET request to retrieve public IP
	response, err := http.Get(URL)
	if err != nil {
		err = oopsBuilder.
			Wrapf(err, "Error occurred while sending GET request to '%s'", URL)

		return "", err
	}
	defer response.Body.Close()

	// Read the response body
	bytesResponse, err := io.ReadAll(response.Body)
	if err != nil {
		err = oopsBuilder.
			With("responseStatusCode", response.StatusCode).
			Wrapf(err, "Error occurred while reading response body")

		return "", err
	}

	// Convert the response body to a string
	return string(bytesResponse), nil
}
