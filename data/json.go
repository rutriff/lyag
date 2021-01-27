package data

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func readBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func FetchJson(token string, url string, proxyUrl *url.URL) ([]byte, error) {
	client := makeClient(proxyUrl)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return readBody(resp)
}

func makeClient(proxyUrl *url.URL) *http.Client {
	var t http.Transport
	if proxyUrl != nil {
		//os.Setenv("HTTP_PROXY", proxyUrl.String())
		t = http.Transport{Proxy: http.ProxyURL(proxyUrl)}
	} else {
		t = http.Transport{}
	}

	return &http.Client{Transport: &t}
}

func PostJson(token string, endpoint string, data []byte, proxyUrl *url.URL) (int, []byte, error) {
	r := bytes.NewReader(data)

	req, err := http.NewRequest("POST", endpoint, r)

	if err != nil {
		return 0, nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "ru-RU")

	client := makeClient(proxyUrl)

	resp, err := client.Do(req)

	if err != nil {
		return 0, nil, err
	}

	body, err := readBody(resp)

	return resp.StatusCode, body, nil
}
