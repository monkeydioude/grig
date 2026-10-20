package compose

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

var ErrNoServices = errors.New("no services mapping in the compose file")

// servicesNode returns the root mapping of a compose file, and the index
// of its "services" key in it.
func servicesNode(content []byte) (*yaml.Node, int, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(content, &doc); err != nil {
		return nil, 0, err
	}
	if len(doc.Content) == 0 || doc.Content[0].Kind != yaml.MappingNode {
		return nil, 0, ErrNoServices
	}
	root := doc.Content[0]
	for i := 0; i+1 < len(root.Content); i += 2 {
		if root.Content[i].Value == "services" && root.Content[i+1].Kind == yaml.MappingNode {
			return root, i, nil
		}
	}
	return nil, 0, ErrNoServices
}

// ServiceNames lists the services of a compose file.
func ServiceNames(content []byte) ([]string, error) {
	root, i, err := servicesNode(content)
	if err != nil {
		return nil, err
	}
	services := root.Content[i+1].Content
	names := make([]string, 0, len(services)/2)
	for j := 0; j < len(services); j += 2 {
		names = append(names, services[j].Value)
	}
	return names, nil
}

// Insert adds svc at the end of the services mapping of content, as text:
// the rest of the file, comments included, is left as it is.
func Insert(content []byte, svc Service) ([]byte, error) {
	root, i, err := servicesNode(content)
	if err != nil {
		return nil, err
	}
	services := root.Content[i+1]
	unit := 2
	if len(services.Content) > 0 {
		unit = services.Content[0].Column - 1
	}
	lines := strings.Split(strings.TrimSuffix(string(content), "\n"), "\n")
	at := insertLine(root, i, lines)
	block := strings.Split(strings.TrimSuffix("\n"+svc.YAML(unit), "\n"), "\n")
	result := append(append(append([]string{}, lines[:at]...), block...), lines[at:]...)
	return []byte(strings.Join(result, "\n") + "\n"), nil
}

// insertLine is the index of the line the services block ends before: the
// next top-level key, or the end of the file. Blank lines and top-level
// comments right before that key stay with it.
func insertLine(root *yaml.Node, servicesKey int, lines []string) int {
	next := servicesKey + 2
	if next >= len(root.Content) {
		return trimTrailingBlank(lines, len(lines))
	}
	at := root.Content[next].Line - 1
	for at > 0 && (strings.TrimSpace(lines[at-1]) == "" || strings.HasPrefix(lines[at-1], "#")) {
		at--
	}
	return at
}

func trimTrailingBlank(lines []string, at int) int {
	for at > 0 && strings.TrimSpace(lines[at-1]) == "" {
		at--
	}
	return at
}

// checkInserted makes sure the file still parses and holds the service.
func checkInserted(content []byte, name string, before int) error {
	names, err := ServiceNames(content)
	if err != nil {
		return fmt.Errorf("compose file broken by the insertion: %w", err)
	}
	if len(names) != before+1 || names[len(names)-1] != name {
		return fmt.Errorf("compose file: %s not added as expected", name)
	}
	return nil
}
