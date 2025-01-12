package model

import (
	"fmt"
	"monkeydioude/grig/internal/errors"
)

type Hook struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Path   string `json:"path"`
	Secret string `json:"secret"`
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
