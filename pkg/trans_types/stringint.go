package trans_types

import (
	"fmt"
	"strconv"

	"gopkg.in/yaml.v3"
)

type StringInt int

// UnmarshalJSON implements `encoding/json`'s `Unmarshaler` interface
// and provides a way to convert a string into an int as part
// of the unmarshaling process of a json.
func (p *StringInt) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(str) > 1 && str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	// Convert the string to an int
	port, err := strconv.Atoi(str)
	if err != nil {
		return fmt.Errorf("invalid port: %w", err)
	}
	*p = StringInt(port)
	return nil
}

// UnmarshalYAML implements `yaml.Unmarshaler` so a port written
// either as an int or as a quoted string decodes into an int.
func (p *StringInt) UnmarshalYAML(node *yaml.Node) error {
	port, err := strconv.Atoi(node.Value)
	if err != nil {
		return fmt.Errorf("invalid port: %w", err)
	}
	*p = StringInt(port)
	return nil
}

func (p StringInt) String() string {
	if p == 0 {
		return ""
	}
	return strconv.Itoa(int(p))
}
