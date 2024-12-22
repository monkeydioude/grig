package model

import (
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

type Capybara struct {
	Proxy    Proxy               `json:"proxy"`
	Services []ServiceDefinition `json:"services"`
}

func (c Capybara) Save() error {
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
