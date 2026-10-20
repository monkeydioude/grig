// Package env_file writes the .env files of services, kept under one base
// directory as <base>/<project>/.env.
package env_file

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	ErrInvalidProject = errors.New("invalid project name: letters, digits, '.', '_' and '-'")
	ErrKeyExists      = errors.New("already set in the .env file")
)

var validProject = regexp.MustCompile(`^[A-Za-z0-9_-][A-Za-z0-9._-]*$`)

// Dir is the base directory holding one folder per project.
type Dir struct {
	Base string
}

// Enabled tells whether a base directory (env_base_path) is configured.
func (d Dir) Enabled() bool {
	return d.Base != ""
}

// Path is the .env file of project, which cannot leave the base directory.
func (d Dir) Path(project string) (string, error) {
	if !validProject.MatchString(project) || strings.Contains(project, "..") {
		return "", ErrInvalidProject
	}
	return filepath.Join(d.Base, project, ".env"), nil
}

// HasKey tells whether the file at path sets key. A missing file sets nothing.
func HasKey(path, key string) (bool, error) {
	content, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return setsKey(content, key), nil
}

// Add appends key=value to the file at path, creating it and its folder,
// readable by their owner only. It refuses a key already set.
func Add(path, key, value string) error {
	has, err := HasKey(path, key)
	if err != nil {
		return err
	}
	if has {
		return fmt.Errorf("%s: %w", key, ErrKeyExists)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()
	prefix, err := separator(path)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, "%s%s=%s\n", prefix, key, value)
	return err
}

// separator is the newline needed before appending, when the file does
// not end with one.
func separator(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	if len(content) == 0 || bytes.HasSuffix(content, []byte("\n")) {
		return "", nil
	}
	return "\n", nil
}

// setsKey tells whether a line of content sets key.
func setsKey(content []byte, key string) bool {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimPrefix(strings.TrimSpace(scanner.Text()), "export ")
		if name, _, ok := strings.Cut(line, "="); ok && strings.TrimSpace(name) == key {
			return true
		}
	}
	return false
}
