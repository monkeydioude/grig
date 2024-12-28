package model

import (
	"encoding/json"
	"fmt"
	"monkeydioude/grig/internal/errors"
	"os"
	"strconv"
)

type Proxy struct {
	Port    int    `json:"port"`
	TLSHost string `json:"tls_host"`
}

type ServiceDefinition struct {
	ID       string `json:"id"`
	Method   string `json:"method"`
	Pattern  string `json:"pattern"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol,omitempty"` // omitempty to handle the absence of this field
}

func (sd ServiceDefinition) PortString() string {
	if sd.Port == 0 {
		return ""
	}
	return strconv.Itoa(sd.Port)
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
		return fmt.Errorf("Capybara.Save(): %w: %w", errors.ErrMarshaling, err)
	}
	if err := c.FileWriter(c.Path, data, os.ModePerm); err != nil {
		return fmt.Errorf("Capybara.Save(): %w: %w", errors.ErrWritingFile, err)
	}
	return nil
}

func (c Capybara) Source() *os.File {
	return nil
}

func (c Capybara) PortString() string {
	if c.Proxy.Port == 0 {
		return ""
	}
	return strconv.Itoa(c.Proxy.Port)
}
