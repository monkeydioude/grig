package app_service

import "monkeydioude/grig/internal/model"

func NewServiceFromPath(path string) (model.Service, error) {
	return NewServiceFromIniPath(path)
}
