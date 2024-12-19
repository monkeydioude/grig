package model

import "os"

type Service struct {
	Environments []string
	Exec         string
	Description  string
}

func (s Service) Save() error {
	return nil
}

func (s Service) Source() *os.File {
	return nil
}
