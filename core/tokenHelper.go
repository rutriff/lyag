package core

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

func FetchToken(url string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	tokenRegexp := regexp.MustCompile("Bearer ([A-z0-9]+)")

	result := tokenRegexp.FindSubmatch(body)

	if result == nil {
		return "", errors.New("no token found")
	}

	return string(result[1]), nil
}
