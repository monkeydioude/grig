package model

import (
	"fmt"
	"monkeydioude/grig/internal/errors"
	"os"

	"gopkg.in/ini.v1"
)

type ServiceSection struct {
	Environments []string `ini:"Environment,allowshadows"`
	Exec         string   `ini:"ExecStart"`
}

type UnitSection struct {
	Description string `ini:"Description"`
}

type Service struct {
	Service ServiceSection `ini:"Service"`
	Unit    UnitSection    `ini:"Unit"`
	Path    string
	Name    string
	OGPath  string `json:"og_path"`
	IniFile *ini.File
}

func (s Service) Save() error {
	if err := s.HydrateIni(s.IniFile); err != nil {
		return err
	}
	if err := s.IniFile.SaveTo(s.Path); err != nil {
		return err
	}
	return nil
}

func (s Service) Source() *os.File {
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

func (s Service) HydrateIni(cfg *ini.File) error {
	if cfg == nil {
		return errors.ErrNilPointer
	}
	fmt.Printf("HydrateIni %+v\n", s)
	cfg.Section("Unit").Key("Description").SetValue(s.Unit.Description)
	serviceSec := cfg.Section("Service")
	serviceSec.Key("ExecStart").SetValue(s.Service.Exec)
	envK := serviceSec.Key("Environment")
	for _, env := range s.Service.Environments {
		envK.AddShadow(env)
	}
	return cfg.SaveTo(s.Path)
}
