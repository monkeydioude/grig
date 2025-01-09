package model

import (
	customErrors "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/pkg/errors"
	"monkeydioude/grig/pkg/model"
)

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

func (c *Deployment) SetParent(p model.IndexBuilder) {
	c.parent = nil
}

func (c Deployment) GetParent() model.IndexBuilder {
	return nil
}

func (c *Deployment) FillBaseData() {
	c.Branches = make([]Branch, 1)
	c.Branches[0].FillBaseData()
	c.Branches[0].SetParent(c)
}

func (c Deployment) Verify() error {
	if c.Repo == "" {
		return errors.Wrap(customErrors.ErrModelVerifyInvalidValue, "Deployment.Verify: Repo")
	}
	return nil
}

func (c *Deployment) VerifyAndSanitize() error {
	if err := c.Verify(); err != nil {
		return err
	}

	for _, br := range c.Branches {
		if err := br.VerifyAndSanitize(); err != nil {
			return err
		}
	}
	return nil
}
