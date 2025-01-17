package services

import (
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/parser"
	"os"
)

type ServiceFilename struct {
	Filename string `json:"filename"`
}

func (sf *ServiceFilename) TryLoadAndParse() (model.Service, error) {
	_, err := os.Stat(sf.Filename)
	if err != nil {
		return model.Service{}, err
	}

	return parser.IniServiceParser(sf.Filename)
}
