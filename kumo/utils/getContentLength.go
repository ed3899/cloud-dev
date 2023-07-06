package utils

import (
	"log"
	"net/http"
)

func GetContentLength(url string) int64 {
	response, err := http.Head(url)
	if err != nil {
		log.Fatalf("there was an error while attempting to download from '%s': '%#v'", url, err)
	}
	defer response.Body.Close()

	return response.ContentLength
}