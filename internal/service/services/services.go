package services

import (
	"fmt"
	"log/slog"
	"monkeydioude/grig/internal/consts"
	"monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/parser"
	"monkeydioude/grig/pkg/dt"
	pkgErr "monkeydioude/grig/pkg/errors"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type AppServicePaths []string

// getIniFile decides whether we need to create a new ini file,
// or load an existing one, according to its presence or absence
// in the AppServicePaths.
func getIniFile(sp AppServicePaths, srvcPath string) (*ini.File, error) {
	ok, path := dt.AnyFunc(sp, func(path string) bool {
		return path == srvcPath
	})
	if !ok {
		return parser.IniNewFile(srvcPath), nil
	}
	return ini.ShadowLoad(path)
}

// cleanIniFileAndSymlink tries to remove an ini file,
// along with its supposed symlink.
// We purposely do not throw err, but log warning
// in case files do not exist.
func cleanIniFileAndSymlink(path string, logger *slog.Logger) {
	name := filepath.Base(path)
	if err := os.Remove(path); err != nil {
		logger.Warn("old ini file remove", "error", err)
	}
	if err := removeSymlink(consts.SYS_SERVICES_DIR, name); err != nil {
		logger.Warn("symlink remove", "error", err)
	}
}

// handleIniSymlink check if the service changed path.
// If so, we delete the old ini file and its symlink,
// and we create a new symlink.
func handleIniSymlink(srvc model.Service, logger *slog.Logger) error {
	if srvc.Path == srvc.OGPath {
		return nil
	}
	cleanIniFileAndSymlink(srvc.OGPath, logger)
	name := filepath.Base(srvc.Path)
	return os.Symlink(srvc.Path, fmt.Sprintf("%s/%s", consts.SYS_SERVICES_DIR, name))
}

// WithUpdate will return a new AppServicePaths
// after being updated with respect to the map
// of services passed.
func (sp AppServicePaths) WithUpdate(
	srvcs *map[string]model.Service,
	logger *slog.Logger,
) (AppServicePaths, error) {
	res := []string{}
	if logger == nil {
		return res, pkgErr.Wrap(errors.ErrNilPointer, "WithUpdate")
	}
	for _, srvc := range *srvcs {
		file, err := getIniFile(sp, srvc.Path)
		if err != nil {
			return res, pkgErr.Wrap(err, "WithUpdate::getIniFile")
		}
		srvc.IniFile = file
		if err := srvc.Save(); err != nil {
			return res, pkgErr.Wrap(err, "WithUpdate::Service.Save")
		}
		if err := handleIniSymlink(srvc, logger); err != nil {
			return res, pkgErr.Wrap(err, "WithUpdate::handleIniSymlink")
		}
		res = append(res, srvc.Path)
	}
	return res, nil
}
