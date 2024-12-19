package app_service

import (
	"fmt"
	"log"
	customErrors "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/model"

	"gopkg.in/ini.v1"
)

func fetchSectionAndKey(cfg *ini.File, section, key string) string {
	unit := cfg.Section(section)
	if unit == nil {
		log.Println("fetchSectionAndKey: invalid section")
		return ""
	}

	sectionKey := unit.Key(key)
	if sectionKey == nil {
		log.Println("fetchSectionAndKey: invalid key")
		return ""
	}

	return sectionKey.String()
}

func fetchSectionAndKeys(cfg *ini.File, section, key string) []string {
	sec := cfg.Section(section)
	if sec == nil || !sec.HasKey(key) {
		log.Println("fetchSectionAndKey: invalid section")
		return []string{}
	}

	sectionKey := sec.Key(key)
	if sectionKey == nil {
		log.Println("fetchSectionAndKey: invalid key")
		return []string{}
	}

	return sectionKey.ValueWithShadows()
}

func NewServiceFromIniPath(path string) (model.Service, error) {
	cfg, err := ini.LoadSources(ini.LoadOptions{
		AllowShadows: true,
	}, path)
	service := model.Service{}
	if err != nil {
		return service, fmt.Errorf("fs.NewServiceFromPath: ini.Load: %w: %w", customErrors.ErrReadIniFile, err)
	}

	service.Description = fetchSectionAndKey(cfg, "Unit", "Description")
	service.Exec = fetchSectionAndKey(cfg, "Service", "ExecStart")
	service.Environments = fetchSectionAndKeys(cfg, "Service", "Environment")
	return service, nil
}
