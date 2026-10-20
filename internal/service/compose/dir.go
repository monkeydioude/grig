package compose

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

var (
	ErrUnknownType = errors.New("unknown compose file type")
	ErrNameTaken   = errors.New("service name already used in")
)

const filePrefix = "docker-compose."

// Dir is the directory holding the docker-compose.<type>.yml files.
// They share one compose project (the directory name), so a service name
// must be unique across all of them.
type Dir struct {
	Base string
}

// Enabled tells whether a directory (docker_compose_base_path) is configured.
func (d Dir) Enabled() bool {
	return d.Base != ""
}

// Files maps each type to its file: docker-compose.tools.yml is "tools".
func (d Dir) Files() (map[string]string, error) {
	entries, err := os.ReadDir(d.Base)
	if err != nil {
		return nil, err
	}
	files := map[string]string{}
	for _, entry := range entries {
		if typ, ok := fileType(entry.Name()); ok && !entry.IsDir() {
			files[typ] = filepath.Join(d.Base, entry.Name())
		}
	}
	return files, nil
}

// Types lists the file types, sorted.
func (d Dir) Types() ([]string, error) {
	files, err := d.Files()
	if err != nil {
		return nil, err
	}
	types := make([]string, 0, len(files))
	for typ := range files {
		types = append(types, typ)
	}
	slices.Sort(types)
	return types, nil
}

func fileType(name string) (string, bool) {
	for _, ext := range []string{".yml", ".yaml"} {
		if typ, ok := strings.CutSuffix(strings.TrimPrefix(name, filePrefix), ext); ok && strings.HasPrefix(name, filePrefix) && typ != "" {
			return typ, true
		}
	}
	return "", false
}

// Add inserts svc in the file of type typ and returns that file's path.
func (d Dir) Add(typ string, svc Service) (string, error) {
	if err := svc.Validate(); err != nil {
		return "", err
	}
	files, err := d.Files()
	if err != nil {
		return "", err
	}
	path, ok := files[typ]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrUnknownType, typ)
	}
	if err := checkNameFree(files, svc.Name); err != nil {
		return "", err
	}
	return path, insertInFile(path, svc)
}

func checkNameFree(files map[string]string, name string) error {
	for _, path := range files {
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		names, err := ServiceNames(content)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if slices.Contains(names, name) {
			return fmt.Errorf("%s: %w %s", name, ErrNameTaken, filepath.Base(path))
		}
	}
	return nil
}

func insertInFile(path string, svc Service) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	before, err := ServiceNames(content)
	if err != nil {
		return err
	}
	updated, err := Insert(content, svc)
	if err != nil {
		return err
	}
	if err := checkInserted(updated, svc.Name, len(before)); err != nil {
		return err
	}
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	return os.WriteFile(path, updated, info.Mode().Perm())
}
