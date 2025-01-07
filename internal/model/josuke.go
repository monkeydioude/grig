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

type Indexer struct {
	Index int `json:"-"`
}

func (i Indexer) GetIndex() int {
	return i.Index
}

func (i *Indexer) SetIndex(index int) {
	i.Index = index
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
	Indexer
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

func (c *Command) InitParent() {
	act := &Action{}
	act.InitParent()
	c.SetParent(act)
}

type Action struct {
	Action   string    `json:"action"`
	Commands []Command `json:"commands"`
	parent   IndexBuilder
	Indexer
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

func (c *Action) FillBaseData() {
	c.Commands = make([]Command, 1)
	c.Commands[0].SetParent(c)
}

func (c *Action) InitParent() {
	br := &Branch{}
	br.InitParent()
	c.SetParent(br)
}

func NewAction(index int) *Action {
	act := Action{}
	act.SetIndex(index)
	act.FillBaseData()
	act.InitParent()
	return &act
}

type Branch struct {
	Branch  string   `json:"branch"`
	Actions []Action `json:"actions"`
	parent  IndexBuilder
	Indexer
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

func (c *Branch) FillBaseData() {
	c.Actions = make([]Action, 1)
	c.Actions[0].FillBaseData()
	c.Actions[0].SetParent(c)
}

func (c *Branch) InitParent() {
	c.SetParent((&Deployment{}))
}

func NewBranch(index int) *Branch {
	br := Branch{}
	br.SetIndex(index)
	br.FillBaseData()
	br.InitParent()
	return &br
}

type Deployment struct {
	Repo     string   `json:"repo"`
	ProjDir  string   `json:"proj_dir"`
	BaseDir  string   `json:"base_dir"`
	Branches []Branch `json:"branches"`
	parent   any
	Indexer
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

func (c *Deployment) FillBaseData() {
	c.Branches = make([]Branch, 1)
	c.Branches[0].FillBaseData()
	c.Branches[0].SetParent(c)
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
	}
}
