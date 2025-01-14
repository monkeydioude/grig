package model

import (
	"errors"
	customErrors "monkeydioude/grig/internal/errors"
	pkgErrors "monkeydioude/grig/pkg/errors"
	"monkeydioude/grig/pkg/model"
	"slices"
)

type Branch struct {
	Branch  string   `json:"branch"`
	Actions []Action `json:"actions"`
	parent  model.IndexBuilder
	Indexer
}

func (Branch) GetName() string {
	return "branches"
}

func (c *Branch) SetParent(p model.IndexBuilder) {
	c.parent = p
}

func (c Branch) GetParent() model.IndexBuilder {
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

func (c Branch) Verify() error {
	if c.Branch == "" {
		return pkgErrors.Wrap(customErrors.ErrModelVerifyInvalidValue, "Branch.Verify: Branch")
	}
	return nil
}

func (c *Branch) VerifyAndSanitize() error {
	if err := c.Verify(); err != nil {
		return err
	}
	if len(c.Actions) == 0 {
		return customErrors.ErrEmptySlice
	}
	for it, act := range slices.Backward(c.Actions) {
		if err := act.VerifyAndSanitize(); err != nil {
			if errors.Is(err, customErrors.ErrModelVerifyInvalidValue) {
				return err
			}
			c.Actions = slices.Delete(c.Actions, it, it+1)
		}
	}
	return nil
}
