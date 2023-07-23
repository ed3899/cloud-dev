package utils

import (
	"net/http"

	"github.com/pkg/errors"
)

func GetContentLength(url string) (int64, error) {
	response, err := http.Head(url)
	if err != nil {
		err = errors.Wrap(err, "failed to get head response")
		return 0, err
	}
	defer response.Body.Close()

	return response.ContentLength, nil
}
