package url

import (
	"net/http"

	"github.com/samber/oops"
)

func GetContentLength(
	url string,
) (
	int64,
	error,
) {
	oopsBuilder :=
		oops.Code("get_content_length_failed").
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
