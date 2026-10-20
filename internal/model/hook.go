package model

import (
	"fmt"
	"monkeydioude/grig/internal/errors"
)

type Hook struct {
	Name   string `json:"name" yaml:"name"`
	Type   string `json:"type" yaml:"type"`
	Path   string `json:"path" yaml:"path"`
	Secret string `json:"secret" yaml:"secret"`
}

func (Hook) GetName() string {
	return "hook"
}

func (h Hook) Verify() error {
	if h.Name == "" {
		return fmt.Errorf("Hook.Verify.Name: %w", errors.ErrModelVerifyInvalidValue)
	}
	if h.Type == "" {
		return fmt.Errorf("Hook.Verify.Type: %w", errors.ErrModelVerifyInvalidValue)
	}
	if h.Path == "" {
		return fmt.Errorf("Hook.Verify.Path: %w", errors.ErrModelVerifyInvalidValue)
	}
	return nil
}
