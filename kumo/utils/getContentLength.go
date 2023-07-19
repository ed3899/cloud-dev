package utils

import (
	"net/http"

	"github.com/pkg/errors"
)

// Returns the content length of the url.
func GetContentLength(url string) (int64, error) {
	// Send a HEAD request to the specified URL
	response, err := http.Head(url)
	if err != nil {
		err = errors.Wrap(err, "failed to get current directory")
		return 0, err
	}
	defer response.Body.Close()

	return response.ContentLength, nil
}
