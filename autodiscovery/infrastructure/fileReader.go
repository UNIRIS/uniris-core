package infrastructure

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileReader struct{}

//ReadFile gets bytes from a file
func (r FileReader) ReadFile(uri string) ([]byte, error) {
	path, err := filepath.Abs(uri)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
