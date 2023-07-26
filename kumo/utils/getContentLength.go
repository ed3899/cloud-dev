package utils

import (
	"net/http"

	"github.com/pkg/errors"
)

func GetContentLength(url string) (contentLength int64, err error) {
	var (
		response *http.Response
	)

	if response, err = http.Head(url); err != nil {
		return 0, errors.Wrap(err, "failed to get head response")
	}
	defer func() {
		if errClosingBody := response.Body.Close(); err != nil {
			err = errors.Wrap(errClosingBody, "failed to close response body")
		}
	}()

	contentLength = response.ContentLength

	return
}
