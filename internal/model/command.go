package model

import (
	"encoding/json"
	"fmt"
	"monkeydioude/grig/pkg/model"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

// Command is the representation of a command from the json config file,
// shaped into a struct, so it gets along well with our chain of Josuke structs.
type Command struct {
	Parts   []string `json:"command"`
	parent  model.IndexBuilder
	Indexer `yaml:"-"`
}

func (Command) GetName() string {
	return "commands"
}

// FlatCommand is the reprensation of a Command in the HTML.
// Instead of Command being a []string, in the html,
// it is, in fact, a json object with a string inside.
// TODO: investigate this
type FlatCommand struct {
	Command string `json:"command"`
}

func (c *Command) UnmarshalJSON(data []byte) error {
	parts := []string{}
	flatCmd := FlatCommand{}
	// not perfect, but we handle both json from the html (flatCmd, one command => "cd /tmp")
	// and from the json file ([]string{}, one command => ["cd", "/tmp"]).

	// First we handle html input...
	if errHtml := json.Unmarshal(data, &flatCmd); errHtml != nil {
		// ...then we try to handle file source
		if errFile := json.Unmarshal(data, &parts); errFile != nil {
			return fmt.Errorf("UnmarshalJSON: %w: %w", errHtml, errFile)
		}
		// That means we handled a file source
		c.Parts = parts
	} else {
		// We handled html input
		c.Parts = strings.Split(flatCmd.Command, " ")
	}
	return nil
}

func (c *Command) MarshalJSON() ([]byte, error) {
	res := make([]string, 0)

	for _, cmd := range c.Parts {
		res = append(res, cmd)
	}
	return json.Marshal(res)
}

// UnmarshalYAML reads a command written as a sequence: ["cd", "/tmp"].
func (c *Command) UnmarshalYAML(node *yaml.Node) error {
	parts := []string{}
	if err := node.Decode(&parts); err != nil {
		return fmt.Errorf("UnmarshalYAML: %w", err)
	}
	c.Parts = parts
	return nil
}

// MarshalYAML writes the command back on one line, ["cd", "/tmp"],
// the way josuke configs are written by hand.
func (c Command) MarshalYAML() (any, error) {
	node := &yaml.Node{Kind: yaml.SequenceNode, Style: yaml.FlowStyle}
	for _, part := range c.Parts {
		node.Content = append(node.Content, &yaml.Node{
			Kind:  yaml.ScalarNode,
			Style: yaml.DoubleQuotedStyle,
			Value: part,
		})
	}
	return node, nil
}

func (c *Command) SetParent(p model.IndexBuilder) {
	c.parent = p
}

func (c Command) GetParent() model.IndexBuilder {
	return c.parent
}

func (c *Command) InitParent() {
	act := &Action{}
	act.InitParent()
	c.SetParent(act)
}

func (c *Command) VerifyAndSanitize() error {
	c.Parts = slices.DeleteFunc(c.Parts, func(cmd string) bool {
		return cmd == ""
	})
	return nil
}

func NewCommand(index int) *Command {
	cmd := Command{}
	cmd.SetIndex(index)
	cmd.InitParent()
	return &cmd
}
