package helpers

import (
	"crypto/md5"
	"io/ioutil"
)

// FileHash returns the hash of a file.
func FileHash(file string) ([16]byte, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return [16]byte{}, err
	}

	return md5.Sum(buf), nil
}
