package model

import "os"

type Josuke struct {
}

func (j Josuke) Save() error {
	return nil
}

func (j Josuke) Source() *os.File {
	return nil
}
