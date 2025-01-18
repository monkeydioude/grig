package model

import (
	"fmt"
	"os"
)

type Service struct {
	Environments []string
	Exec         string
	Description  string
	Path         string
	Name         string
}

func (s Service) Save() error {
	// @Todo: save ini file
	return nil
}

func (s Service) Source() *os.File {
	return nil
}

func (s Service) IdGen(key string) string {
	return fmt.Sprintf("%s[%s]", s.Name, key)
}

func (s Service) EnvironmentIdGen(it int) string {
	return fmt.Sprintf("%s[%d]", s.IdGen("environments"), it)
}
