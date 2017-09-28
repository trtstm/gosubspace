package helpers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

// AssertNoError panics on error.
func AssertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

var client = &http.Client{Timeout: 10 * time.Second}

// GetJSON retrieves json from an url.
func GetJSON(url string, target interface{}) error {
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// DownloadFile will download a file from url and put in in filepath.
func DownloadFile(url, filepath string) (err error) {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
