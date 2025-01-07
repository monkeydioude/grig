package model

import (
	"encoding/json"
	"errors"
	"fmt"
	customErrors "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/pkg/trans_types"
	"os"
	"strconv"
)

type Proxy struct {
	Port    trans_types.StringInt `json:"port"`
	TLSHost string                `json:"tls_host"`
}

func (p Proxy) Verify() error {
	if p.Port <= 0 {
		return fmt.Errorf("Proxy.Verify(): Port: %d: %w", p.Port, customErrors.ErrModelVerifyInvalidValue)
	}
	if p.TLSHost == "" {
		return fmt.Errorf("Proxy.Verify(): TLSHost: %q: %w", p.TLSHost, customErrors.ErrModelVerifyInvalidValue)
	}
	return nil
}

type ServiceDefinition struct {
	ID       string                `json:"id"`
	Method   string                `json:"method"`
	Pattern  string                `json:"pattern"`
	Port     trans_types.StringInt `json:"port"`
	Protocol string                `json:"protocol,omitempty"` // omitempty to handle the absence of this field
}

func (sd ServiceDefinition) Verify() error {
	if sd.ID == "" {
		return fmt.Errorf("Proxy.Verify(): ID: %q: %w", sd.ID, customErrors.ErrModelVerifyInvalidValue)
	}
	if sd.Method == "" {
		return fmt.Errorf("Proxy.Verify(): Method: %q: %w", sd.Method, customErrors.ErrModelVerifyInvalidValue)
	}
	if sd.Pattern == "" {
		return fmt.Errorf("Proxy.Verify(): Pattern: %q: %w", sd.Pattern, customErrors.ErrModelVerifyInvalidValue)
	}
	if sd.Port <= 0 {
		return fmt.Errorf("Proxy.Verify(): Port: %d: %w", sd.Port, customErrors.ErrModelVerifyInvalidValue)
	}
	return nil
}

func (sd ServiceDefinition) PortString() string {
	if sd.Port == 0 {
		return ""
	}
	return strconv.Itoa(int(sd.Port))
}

type Capybara struct {
	Proxy      Proxy                                   `json:"proxy"`
	Services   []ServiceDefinition                     `json:"services"`
	Path       string                                  `json:"-"`
	FileWriter func(string, []byte, os.FileMode) error `json:"-"`
}

func (c Capybara) Save() error {
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("Capybara.Save(): %w: %w", customErrors.ErrMarshaling, err)
	}
	if c.FileWriter == nil {
		return fmt.Errorf("Capybara.Save(): c.FileWriter: %w", customErrors.ErrNilPointer)
	}
	if err := c.FileWriter(c.Path, data, os.ModePerm); err != nil {
		return fmt.Errorf("Capybara.Save(): %w: %w", customErrors.ErrWritingFile, err)
	}
	return nil
}

func (c Capybara) Verify() error {
	var errs error
	if err := c.Proxy.Verify(); err != nil {
		errs = errors.Join(errs, err)
	}

	for _, sd := range c.Services {
		if err := sd.Verify(); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (c Capybara) CloneBase() Capybara {
	return Capybara{
		Path:       c.Path,
		FileWriter: c.FileWriter,
	}
}
