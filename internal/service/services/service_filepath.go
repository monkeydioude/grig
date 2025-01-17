package services

import (
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/parser"
	"os"
)

type Filepath struct {
	Filepath string `json:"filepath"`
}

func (sf *Filepath) TryLoadAndParse() (model.Service, error) {
	_, err := os.Stat(sf.Filepath)
	if err != nil {
		return model.Service{}, err
	}

	return parser.IniServiceParser(sf.Filepath)
}
