package parser

import (
	"errors"
	"log/slog"
	customErrors "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/model"
	pkgErr "monkeydioude/grig/pkg/errors"
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

func IniNewFile(path string) *ini.File {
	return ini.Empty(ini.LoadOptions{
		AllowShadows: true,
	})
}

func IniServiceParser(path string) (model.Service, error) {
	service := model.Service{}
	cfg, err := ini.ShadowLoad(path)
	if err != nil {
		return service, pkgErr.Wrapf(err, "fs.IniServiceParser: ini.Load: %w", customErrors.ErrReadIniFile)
	}
	service.Path = path
	service.Name = filepath.Base(path)
	if err := cfg.MapTo(&service); err != nil {
		return service, pkgErr.Wrapf(err, "cfg.MapTo: %w", customErrors.ErrReadIniFile)
	}
	serviceSection := cfg.Section("Service")
	// mapping by hand, cause the pkg doesnt handle it well
	service.Service.Environment = serviceSection.Key("Environment").ValueWithShadows()
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
