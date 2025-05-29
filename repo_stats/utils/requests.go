package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"
)

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

func Get(url string, body string, headers map[string]string) (string, error) {
	resp, err := makeRequest("GET", url, body, headers)
	if err != nil {
		return "", WrapError(err, "get", "while calling makeRequest")
	}

	// Read response
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", WrapError(err, "get", "while reading response")
	}
	respBodyString := string(respBody)
	return respBodyString, nil
}

func ParseBody(body string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(body), &result)
	if err != nil {
		return nil, WrapError(err, "parseBody", "while parsing body")
	}
	return result, nil
}

// ExtractJson
//
// Extract the contents of a json map
// Parameters:
//   - parsedJson: map of string to interface from json.Unmarshal([]byte(body), &result)
//   - key: The key to look for in the json
//
// # Returns T type result found in the json, error on failure
//
// Example usage:
//
//	extractJson[string](parsedBody, "name")
//	   - for a "name" that maps to a string
//	extractJson[map[string]interface{}](parsedBody, "moreData")
//	   - for a "moreData" that maps to another map
func ExtractJson[T any](parsedJson map[string]interface{}, key string) (T, error) {
	var zero T
	value, exists := parsedJson[key]
	if !exists {
		return zero, WrapError(errors.New(key), "extractJson", "key not found")
	}

	result, ok := value.(T)
	if !ok {
		return zero, WrapError(errors.New(key), "extractJson",
			"value is not a \""+reflect.TypeOf(zero).String()+
				"\" but a \""+reflect.TypeOf(value).String()+"\"")
	}

	return result, nil
}
