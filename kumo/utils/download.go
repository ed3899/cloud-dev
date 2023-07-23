package utils

import (
	"net/http"

	"github.com/pkg/errors"
)

func Download(url string, dest string) (err error) {
	response, err := http.Get(url)
	if err != nil {
		err = errors.Wrapf(err, "failed to download from: %s", url)
		return
	}
	defer response.Body.Close()
}