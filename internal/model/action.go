package model

import (
	customErrors "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/pkg/errors"
	"monkeydioude/grig/pkg/model"
	"slices"
)

type Action struct {
	Action   string    `json:"action"`
	Commands []Command `json:"commands"`
	parent   model.IndexBuilder
	Indexer
}

func (Action) GetName() string {
	return "actions"
}

func (c *Action) SetParent(p model.IndexBuilder) {
	c.parent = p
}

func (c Action) GetParent() model.IndexBuilder {
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

func (c Action) Verify() error {
	if c.Action == "" {
		return errors.Wrap(customErrors.ErrModelVerifyInvalidValue, "Action.Verify: Action")
	}
	return nil
}

func (c *Action) VerifyAndSanitize() error {
	if err := c.Verify(); err != nil {
		return err
	}

	c.Commands = slices.DeleteFunc(c.Commands, func(cmd Command) bool {
		return cmd.VerifyAndSanitize() != nil
	})
	return nil
}
