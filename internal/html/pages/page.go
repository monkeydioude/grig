package pages

import (
	"github.com/a-h/templ"
)

type Page interface {
	Title() string
	Content() templ.Component
}
