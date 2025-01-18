package parser

import (
	"errors"
	"fmt"
	"log/slog"
	customErrors "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/model"
	"path/filepath"

	"gopkg.in/ini.v1"
)

func fetchSectionAndKey(cfg *ini.File, section, key string) string {
	unit := cfg.Section(section)
	if unit == nil {
		slog.Info("fetchSectionAndKey: invalid section")
		return ""
	}

	sectionKey := unit.Key(key)
	if sectionKey == nil {
		slog.Info("fetchSectionAndKey: invalid key")
		return ""
	}

	return sectionKey.String()
}

func fetchSectionAndKeys(cfg *ini.File, section, key string) []string {
	sec := cfg.Section(section)
	if sec == nil || !sec.HasKey(key) {
		slog.Info("fetchSectionAndKey: invalid section")
		return []string{}
	}

	sectionKey := sec.Key(key)
	if sectionKey == nil {
		slog.Info("fetchSectionAndKey: invalid key")
		return []string{}
	}

	return sectionKey.ValueWithShadows()
}

func IniServiceParser(path string) (model.Service, error) {
	cfg, err := ini.LoadSources(ini.LoadOptions{
		AllowShadows: true,
	}, path)
	service := model.Service{}
	if err != nil {
		return service, fmt.Errorf("fs.NewServiceFromPath: ini.Load: %w: %w", customErrors.ErrReadIniFile, err)
	}
	service.Path = path
	service.Name = filepath.Base(path)
	service.Description = fetchSectionAndKey(cfg, "Unit", "Description")
	service.Exec = fetchSectionAndKey(cfg, "Service", "ExecStart")
	service.Environments = fetchSectionAndKeys(cfg, "Service", "Environment")
	return service, nil
}

func IniServicesParser(paths []string) ([]model.Service, error) {
	res := make([]model.Service, 0, len(paths))
	var errs error
	for _, path := range paths {
		srvc, err := IniServiceParser(path)
		if err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		res = append(res, srvc)
	}
	return res, errs
}
