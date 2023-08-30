package url

import (
	"net/http"

	"github.com/samber/oops"
)

// Returns the content length of the given URL.
//
// Example:
//
//	("https://releases.hashicorp.com/packer/1.7.4/packer_1.7.4_windows_amd64.zip") -> 123456789, nil
func GetContentLength(
	url string,
) (
	int64,
	error,
) {
	oopsBuilder :=
		oops.Code("GetContentLength").
			With("url", url)

	response, err := http.Head(url)
	if err != nil {
		err := oopsBuilder.
			Wrapf(err, "failed to get head response from: %s", url)
		return 0, err
	}
	defer response.Body.Close()

	return response.ContentLength, nil
}
