package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var (
	// Client holds the HTTP client object that is used by the HTTP utils to make HTTP calls.
	// By default, the http.DefaultClient is used, but this can be changed if needed.
	Client = http.DefaultClient
)

// GetJSON is a function that reads a JSON document from a given URL and marshals that into a given result object.
func GetJSON(url string, result interface{}) error {
	response, err := Client.Get(url) //nolint:gosec
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return fmt.Errorf("unexpected status while retrieving releases: %s", response.Status)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return err
	}

	return nil
}

// GetFile downloads a given URL into a destination file.
// If the flag overwrite is set to false, the destination file will not be overwritten and nothing will be downloaded.
func GetFile(url, destinationFile string, overwrite bool) (bool, error) {
	if PathExists(destinationFile) && !overwrite {
		return false, nil
	}

	TryRemove(destinationFile)

	response, err := Client.Get(url) //nolint:gosec
	if err != nil {
		return true, err
	}
	if response.StatusCode != 200 {
		return true, fmt.Errorf("unexpected status while retrieving release file: %s", response.Status)
	}

	defer func() {
		_ = response.Body.Close()
	}()

	directory, _ := filepath.Split(destinationFile)
	if err := os.MkdirAll(directory, 0755); err != nil {
		return true, err
	}

	file, err := os.Create(destinationFile)
	if err != nil {
		return true, err
	}

	defer func() {
		_ = file.Close()
	}()

	if _, err := io.Copy(file, response.Body); err != nil {
		return true, err
	}

	return true, nil
}

// StaticResponseClient is a function that create a HTTP client that always produces the same response.
// This is primarily used by tests.
func StaticResponseClient(statusCode int, body []byte, err error) *http.Client {
	return &http.Client{
		Transport: staticResponseRoundTripper{
			statusCode: statusCode,
			body:       body,
			err:        err,
		},
		CheckRedirect: http.DefaultClient.CheckRedirect,
		Jar:           http.DefaultClient.Jar,
		Timeout:       http.DefaultClient.Timeout,
	}
}

type staticResponseRoundTripper struct {
	statusCode int
	body       []byte
	err        error
}

func (rt staticResponseRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	if rt.err != nil {
		return nil, rt.err
	}

	resp := &http.Response{
		StatusCode: rt.statusCode,
		Header:     make(http.Header),
		Body:       ioutil.NopCloser(bytes.NewBuffer(rt.body)),
	}

	return resp, nil
}
