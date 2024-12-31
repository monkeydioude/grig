package pages

import "monkeydioude/grig/internal/html/element"

type Error string

func (e Error) Title() string {
	return string(e)
}

func (Error) Link() element.Link {
	return element.Link{}
}
