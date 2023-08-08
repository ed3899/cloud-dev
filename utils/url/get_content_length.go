package url

import (
	"log"
	"net/http"

	"github.com/samber/oops"
)

func GetContentLength(
	url string,
) (
	contentLength int64,
	err error,
) {
	var (
		oopsBuilder = oops.Code("get_content_length_failed").
				With("url", url)

		response *http.Response
	)

	if response, err = http.Head(url); err != nil {
		err = oopsBuilder.
			Wrapf(err, "failed to get head response from: %s", url)
		return
	}
	defer func(response *http.Response) {
		if err := response.Body.Close(); err != nil {
			log.Fatalf(
				"%+v",
				oopsBuilder.
					With("responseStatus", response.Status).
					Wrapf(err, "failed to close response body"),
			)
		}
	}(response)

	contentLength = response.ContentLength

	return
}

type GetContentLengthF func(
	url string,
) (
	contentLength int64,
	err error,
)
