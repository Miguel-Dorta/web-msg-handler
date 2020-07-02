package client
// Package client represents the HTTP client that is used internally by web-msg-handler to make HTTP request.

import (
	"bytes"
	"fmt"
	"github.com/Miguel-Dorta/web-msg-handler/pkg/mime"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// c is the common HTTP client for this package
var c = &http.Client{Timeout: 10 * time.Second}

//TODO decide if remove PostJSON. Not longer in use.

// PostJSON makes a POST request to the URL provided with the JSON data provided.
// It returns the data of the body of the response.
func PostJSON(url string, data []byte) ([]byte, error) {
	return processResponse(c.Post(url, mime.JSON, bytes.NewReader(data)))
}

// PostJSON makes a POST request to the URL provided with the url form data provided.
// It returns the data of the body of the response.
func PostForm(url string, data url.Values) ([]byte, error) {
	return processResponse(c.PostForm(url, data))
}

func processResponse(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, fmt.Errorf("http request failed: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)
	}

	if err = resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("error closing response body: %s", err)
	}

	return body, nil
}

