package htmlinput

import (
	"monkeydioude/grig/pkg/model"
	"regexp"
	"slices"
	"strconv"
)

func FillIndexesFromHTMLInput(hp model.IndexBuilder, htmlInput string) error {
	if hp == nil {
		return nil
	}
	// Capture numbers surrounded by brackets only
	pattern := `\[(\d+)\]`
	// Find all matches
	matches := regexp.MustCompile(pattern).FindAllStringSubmatch(htmlInput, -1)
	// Starting from the end
	slices.Reverse(matches)
	// Extract matched groups
	for _, match := range matches {
		// Check capturing group
		if match[1] != "" {
			index, err := strconv.Atoi(match[1])
			if err != nil {
				return err
			}
			hp.SetIndex(index)
		}
		if hp.GetParent() == nil {
			return nil
		}
		hp = hp.GetParent()
	}
	return nil
}
