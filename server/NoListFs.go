package server

import (
	"net/http"
	"fmt"
)

type NoListFs string

func (n NoListFs) Open(name string) (http.File, error) {
	f, err := http.Dir(name).Open(name)
	if err != nil {
		return nil, err
	}

	if i, err := f.Stat(); err != nil {
		return nil, fmt.Errorf("error with file stat: %v", err)
	} else if i.IsDir() {
		return nil, fmt.Errorf("error opening file: is a directory")
	}

	return f, nil
}