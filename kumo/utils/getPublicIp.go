package utils

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func GetPublicIp() (ip string, err error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		err = errors.Wrap(err, "Error occurred while getting public IP")
		return "", err
	}
	defer resp.Body.Close()

	bytesResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "Error occurred while reading response body")
		return "", err
	}

	ip = string(bytesResp)

	return ip, nil
}
