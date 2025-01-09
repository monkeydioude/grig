package model

import (
	"encoding/json"
	"fmt"
	"log"
	customErr "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/pkg/errors"
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
		return fmt.Errorf("Proxy.Verify(): Port: %d: %w", p.Port, customErr.ErrModelVerifyInvalidValue)
	}
	if p.TLSHost == "" {
		return fmt.Errorf("Proxy.Verify(): TLSHost: %q: %w", p.TLSHost, customErr.ErrModelVerifyInvalidValue)
	}
	return nil
}

type ServiceDefinition struct {
	ID       string                `json:"id"`
	Method   string                `json:"method"`
	Pattern  string                `json:"pattern"`
	Port     trans_types.StringInt `json:"port"`
	Protocol string                `json:"protocol,omitempty"`
}

func wrapErr(method, id string) error {
	return errors.Wrapf(customErr.ErrModelVerifyInvalidValue, "ServiceDefinition.%s(): ID: %q", method, id)
}

func (sd ServiceDefinition) Verify() error {
	if sd.ID == "" {
		return wrapErr("Verify", sd.ID)
	}
	if sd.Method == "" {
		return wrapErr("Method", sd.ID)
	}
	if sd.Pattern == "" {
		return wrapErr("Pattern", sd.ID)
	}
	if sd.Port <= 0 {
		return wrapErr("Port", sd.ID)
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
		return errors.Wrapf(err, "Capybara.Save(): %w", customErr.ErrMarshaling)
	}
	if c.FileWriter == nil {
		return errors.Wrapf(customErr.ErrNilPointer, "Capybara.FileWriter()")
	}
	if err := c.FileWriter(c.Path, data, os.ModePerm); err != nil {
		return errors.Wrapf(err, "Capybara.Save(): %w", customErr.ErrWritingFile)
	}
	return nil
}

func (c *Capybara) Sanitize() {
	// c.Services = slices.DeleteFunc(c.Services, func (sd ServiceDefinition) bool {
	// 	return sd.Verify() != nil
	// })
	// use of a custom index, so we dont range into a non existant element
	// in case we delete one
	j := 0
	for i := 0; i < len(c.Services); i++ {
		sd := c.Services[i]
		if err := sd.Verify(); err != nil {
			log.Printf("[ERR ] Capybara.Sanitize: %+v", err.Error())
			// Skip the element by not copying it to the new position
			continue
		}
		c.Services[j] = c.Services[i]
		j++
	}
	c.Services = c.Services[:j]
}

func (c Capybara) CloneBase() Capybara {
	return Capybara{
		Path:       c.Path,
		FileWriter: c.FileWriter,
	}
}
