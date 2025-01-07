package model

import (
	"encoding/json"
	"fmt"
	"monkeydioude/grig/pkg/trans_types"
	"os"
)

type IndexBuilder interface {
	SetIndex(int)
	GetIndex() int
	GetParent() IndexBuilder
	SetParent(IndexBuilder)
	GetName() string
}

type indexer struct {
	index int `json:"-"`
}

func (i indexer) GetIndex() int {
	return i.index
}

func (i *indexer) SetIndex(index int) {
	i.index = index
}

type Hook struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Path   string `json:"path"`
	Secret string `json:"secret"`
}

func (Hook) GetName() string {
	return "hook"
}

type Command struct {
	Command []string `json:"commands"`
	parent  IndexBuilder
	indexer
}

func (Command) GetName() string {
	return "command"
}

func (c *Command) UnmarshalJSON(data []byte) error {
	cmds := []string{}

	err := json.Unmarshal(data, &cmds)
	if err != nil {
		return fmt.Errorf("UnmarshalJSON: %w", err)
	}
	c.Command = cmds
	return nil
}

func (c *Command) SetParent(p IndexBuilder) {
	c.parent = p
}

func (c Command) GetParent() IndexBuilder {
	return c.parent
}

type Action struct {
	Action   string    `json:"action"`
	Commands []Command `json:"commands"`
	parent   IndexBuilder
	indexer
}

func (Action) GetName() string {
	return "action"
}

func (c *Action) SetParent(p IndexBuilder) {
	c.parent = p
}

func (c Action) GetParent() IndexBuilder {
	return c.parent
}

type Branch struct {
	Branch  string   `json:"branch"`
	Actions []Action `json:"actions"`
	parent  IndexBuilder
	indexer
}

func (Branch) GetName() string {
	return "branch"
}

func (c *Branch) SetParent(p IndexBuilder) {
	c.parent = p
}

func (c Branch) GetParent() IndexBuilder {
	return c.parent
}

type Deployment struct {
	Repo     string   `json:"repo"`
	ProjDir  string   `json:"proj_dir"`
	BaseDir  string   `json:"base_dir"`
	Branches []Branch `json:"branches"`
	parent   any
	indexer
}

func (Deployment) GetName() string {
	return "deployment"
}

func (c *Deployment) SetParent(p IndexBuilder) {
	c.parent = nil
}

func (c Deployment) GetParent() IndexBuilder {
	return nil
}

type Josuke struct {
	LogLevel         string                `json:"logLevel"`
	Host             string                `json:"host"`
	Port             trans_types.StringInt `json:"port"`
	Store            string                `json:"store"`
	HealthcheckRoute string                `json:"healthcheck_route"`
	Hook             []Hook                `json:"hook"`
	Deployment       []Deployment          `json:"deployment"`
}

func (j Josuke) Save() error {
	return nil
}

func (j Josuke) Source() *os.File {
	return nil
}

func (j *Josuke) FillBaseData() {
	if j.HealthcheckRoute == "" {
		j.HealthcheckRoute = "/josuke/healthcheck"
	}
	if len(j.Hook) == 0 {
		j.Hook = make([]Hook, 1)
	}
	if len(j.Deployment) == 0 {
		j.Deployment = make([]Deployment, 1)
		j.Deployment[0].Branches = make([]Branch, 1)
		j.Deployment[0].Branches[0].Actions = make([]Action, 1)
		j.Deployment[0].Branches[0].Actions[0].Commands = make([]Command, 1)
	}
}
