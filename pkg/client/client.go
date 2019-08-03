package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const mimeJSON = "application/json"

var c = http.Client{Timeout: 10 * time.Second}

func PostJSON(url string, data []byte) ([]byte, error) {
	resp, err := c.Post(url, mimeJSON, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed http request: %s", err)
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

