package model

import (
	"fmt"
	"monkeydioude/grig/internal/errors"

	"gopkg.in/ini.v1"
)

type Service struct {
	Service ServiceSection `ini:"Service"`
	Unit    UnitSection    `ini:"Unit"`
	Path    string
	Name    string
	OGPath  string `json:"og_path"`
	IniFile *ini.File
}

func (s Service) Save() error {
	if err := s.hydrateIni(s.IniFile); err != nil {
		return err
	}
	if err := s.IniFile.SaveTo(s.Path); err != nil {
		return err
	}
	return nil
}

func (s Service) IdGen(it int, section, key string) string {
	if section != "" {
		section = "[" + section + "]"
	}
	return fmt.Sprintf("services[%d]%s[%s]", it, section, key)
}

func (s Service) EnvironmentIdGen(base string, it int) string {
	return fmt.Sprintf("%s[environment][%d]", base, it)
}

func (s Service) hydrateIni(cfg *ini.File) error {
	if cfg == nil {
		return errors.ErrNilPointer
	}
	unit := cfg.Section("Unit")
	unit.Key("Description").SetValue(s.Unit.Description)
	unit.Key("After").SetValue(string(s.Unit.After))

	serviceSec := cfg.Section("Service")
	serviceSec.Key("ExecStart").SetValue(s.Service.ExecStart)
	serviceSec.Key("Type").SetValue(string(s.Service.Type))
	// remake Environment key
	serviceSec.DeleteKey("Environment")
	envK := serviceSec.Key("Environment")
	for _, env := range s.Service.Environment {
		envK.AddShadow(env)
	}

	return nil
}

type ServiceType string

const (
	SimpleService ServiceType = "simple"
)

type ServiceSection struct {
	Environment []string    `ini:"Environment,allowshadow" json:"environment"`
	ExecStart   string      `ini:"ExecStart,allowshadow" json:"exec_start"`
	Type        ServiceType `ini:"Type,allowshadow"`
}

type UnitAfter string

const (
	NetworkOnline UnitAfter = "network-online.target"
)

type UnitSection struct {
	Description string    `ini:"Description"`
	After       UnitAfter `ini:"After"`
}
