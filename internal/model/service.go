package model

import (
	"fmt"
	"monkeydioude/grig/internal/errors"

	"gopkg.in/ini.v1"
)

type Service struct {
	Service ServiceSection
	Unit    UnitSection
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

func (s Service) IdGen(section, key string) string {
	if section != "" {
		section = "[" + section + "]"
	}
	return fmt.Sprintf("%s%s[%s]", s.Name, section, key)
}

func (s Service) EnvironmentIdGen(it int) string {
	return fmt.Sprintf("%s[%d]", s.IdGen("service", "environments"), it)
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
	Environment []string
	ExecStart   string
	Type        ServiceType
}

type UnitAfter string

const (
	NetworkOnline UnitAfter = "network-online.target"
)

type UnitSection struct {
	Description string
	After       UnitAfter
}
