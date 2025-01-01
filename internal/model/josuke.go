package model

import (
	"monkeydioude/grig/internal/service/json_types"
	"os"
)

type Hook struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Path   string `json:"path"`
	Secret string `json:"secret"`
}

type Command struct {
	Command []string `json:"commands"`
}

type Action struct {
	Action   string     `json:"action"`
	Commands [][]string `json:"commands"`
}

type Branch struct {
	Branch  string   `json:"branch"`
	Actions []Action `json:"actions"`
}

type Deployment struct {
	Repo     string   `json:"repo"`
	ProjDir  string   `json:"proj_dir"`
	BaseDir  string   `json:"base_dir"`
	Branches []Branch `json:"branches"`
}

type Josuke struct {
	LogLevel         string               `json:"logLevel"`
	Host             string               `json:"host"`
	Port             json_types.StringInt `json:"port"`
	Store            string               `json:"store"`
	HealthcheckRoute string               `json:"healthcheck_route"`
	Hook             []Hook               `json:"hook"`
	Deployment       []Deployment         `json:"deployment"`
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
		j.Deployment[0].Branches[0].Actions[0].Commands = make([][]string, 1)

	}
}
