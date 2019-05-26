package services

import (
	"fmt"
	"io"
	"io/ioutil"
)

type Service struct {
	PasswdPath     string
	GroupPath      string
	readerBuilders map[string]func(source string) io.Reader
}

// ReadData will return byte data that corresponds to a data source represented by a string value
func (s Service) ReadData(source string) ([]byte, error) {

	var readerBuilder func(source string) io.Reader
	var exists bool

	if readerBuilder, exists = s.readerBuilders[source]; !exists {
		return nil, fmt.Errorf("reader does not exist for source: %s\n", source)
	}

	result, err := ioutil.ReadAll(readerBuilder(source))

	return result, err
}

func NewService(passwdPath, groupPath string,
	readerBuilders map[string]func(source string) io.Reader) *Service {

	return &Service{
		PasswdPath:     passwdPath,
		GroupPath:      groupPath,
		readerBuilders: readerBuilders,
	}
}
