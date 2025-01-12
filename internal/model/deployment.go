package model

import (
	"errors"
	customErrors "monkeydioude/grig/internal/errors"
	pkgErrors "monkeydioude/grig/pkg/errors"
	"monkeydioude/grig/pkg/model"
	"slices"
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
		return pkgErrors.Wrap(customErrors.ErrModelVerifyInvalidValue, "Deployment.Verify: Repo")
	}
	return nil
}

func (c *Deployment) VerifyAndSanitize() error {
	if err := c.Verify(); err != nil {
		return err
	}
	if len(c.Branches) == 0 {
		return customErrors.ErrEmptySlice
	}
	for it, br := range slices.Backward(c.Branches) {
		if err := br.VerifyAndSanitize(); err != nil {
			if errors.Is(err, customErrors.ErrModelVerifyInvalidValue) {
				return err
			}
			c.Branches = slices.Delete(c.Branches, it, it+1)
		}
	}
	return nil
}
