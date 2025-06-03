package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"
)

// makeRequest
// Makes a request using the "net/http" package
//
// Parameters:
//   - method: method for request, such as "get", "put", etc
//   - url: url to send request to
//   - body: body content to send alongside url, use empty string for no body
//   - headers: map of header key value pairs to send with request
//
// Returns response pointer and any errors
func makeRequest(method string, url string, body string, headers map[string]string) (*http.Response, error) {
	// Create client
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Create request
	var readerBody io.Reader = nil
	if body != "" {
		readerBody = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, url, readerBody)
	if err != nil {
		return nil, WrapError(err, "makeRequest", "while creating request")
	}
	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, WrapError(err, "makeRequest", "while sending request")
	}

	// Handle codes
	if resp.StatusCode != http.StatusOK {
		return nil, WrapError(errors.New(resp.Status), "get", "status code error")
	}

	return resp, nil
}

// Get
// Makes get request
//
// Parameters:
//   - url: url to make request to
//   - body: body content to send alongside url, use empty string for no body
//   - headers: map of header key value pairs to send with request
//
// Returns response pointer and any errors
func Get(url string, body string, headers map[string]string) (string, http.Header, error) {
	resp, err := makeRequest("GET", url, body, headers)
	if err != nil {
		return "", nil, WrapError(err, "get", "while calling makeRequest")
	}

	// Read response
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, WrapError(err, "get", "while reading response")
	}
	respBodyString := string(respBody)

	respHeader := resp.Header
	return respBodyString, respHeader, nil
}

// ParseBody
// Parses response body to create JSON interface
//
// Parameters:
//   - body: body content to parse
//
// Returns JSON interface
func ParseBody(body string) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal([]byte(body), &result)
	if err != nil {
		return nil, WrapError(err, "parseBody", "while parsing body")
	}
	return result, nil
}
