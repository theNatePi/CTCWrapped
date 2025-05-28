package main

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

func makeRequest(method string, url string, body string, headers map[string]string) (*http.Response, error) {
	// Create client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Create request
	readerBody := strings.NewReader(body)
	req, err := http.NewRequest(method, url, readerBody)
	if err != nil {
		return nil, wrapError(err, "makeRequest", "while creating request")
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, wrapError(err, "makeRequest", "while sending request")
	}

	// Handle codes
	if resp.StatusCode != http.StatusOK {
		return nil, wrapError(errors.New(resp.Status), "get", "status code error")
	}

	return resp, nil
}

func get(url string, body string, headers map[string]string) ([]byte, error) {
	resp, err := makeRequest("GET", url, body, headers)
	if err != nil {
		return nil, wrapError(err, "get", "while calling makeRequest")
	}

	// Read response
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, wrapError(err, "get", "while reading response")
	}
	return respBody, nil
}
