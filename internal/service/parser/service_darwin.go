package parser

import "monkeydioude/grig/internal/model"

func ServiceFileParser(path string) (model.Service, error) {
	return IniServiceParser(path)
}
