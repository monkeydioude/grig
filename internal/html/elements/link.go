package elements

import (
	"fmt"

	"github.com/a-h/templ"
)

type Target string

const (
	Self   Target = "_self"
	Blank  Target = "_blank"
	Parent Target = "_parent"
)

func (t Target) String() string {
	return string(t)
}

type Link struct {
	Href   templ.SafeURL
	Text   fmt.Stringer
	Target Target
}

type Text string

func (t Text) String() string {
	return string(t)
}

func MainNavigation() []Link {
	return []Link{
		{
			Href:   "/",
			Text:   Text("Home"),
			Target: Self,
		},
	}
}
