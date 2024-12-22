package main

import "monkeydioude/grig/internal/html/element"

func MainNavigation() element.Nav {
	return element.Nav{
		Links: []element.Link{
			{
				Href:   "/",
				Text:   element.Text("Home"),
				Target: element.Self,
			},
			{
				Href:   "/capybara",
				Text:   element.Text("Capybara"),
				Target: element.Self,
			},
		},
	}
}
