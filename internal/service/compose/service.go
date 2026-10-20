// Package compose adds the service of a new project to one of the
// docker-compose.<type>.yml files of a directory.
package compose

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrInvalidName    = errors.New("invalid service name: letters, digits, '_' and '-'")
	ErrInvalidPort    = errors.New("invalid port")
	ErrInvalidContext = errors.New("invalid build context: an absolute path")
	ErrInvalidHealth  = errors.New("invalid healthcheck path: starts with /")
)

var (
	validName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)
	safePath  = regexp.MustCompile(`^/[A-Za-z0-9._/-]*$`)
)

// Service is a compose service built from its repository, on the host
// network like the other services, with its .env and a healthcheck.
type Service struct {
	Name         string
	BuildContext string
	Dockerfile   string
	Port         int
	// HealthPath is checked on localhost:Port. Empty: no healthcheck.
	HealthPath string
	// EnvFile is loaded when it exists. Empty: no env file.
	EnvFile string
}

func (s Service) Validate() error {
	switch {
	case !validName.MatchString(s.Name):
		return ErrInvalidName
	case s.Port <= 0 || s.Port > 65535:
		return ErrInvalidPort
	case !safePath.MatchString(s.BuildContext):
		return ErrInvalidContext
	case s.HealthPath != "" && !safePath.MatchString(s.HealthPath):
		return ErrInvalidHealth
	}
	return nil
}

// YAML renders the service as an entry of a services mapping, indented
// by unit spaces per level, like the file it goes in. Service must be
// valid: values are inlined.
func (s Service) YAML(unit int) string {
	lines := []string{
		s.Name + ":",
		"  build:",
		"    context: " + s.BuildContext,
		"    dockerfile: " + s.Dockerfile,
		"  labels:",
		`    logging: "true"`,
		fmt.Sprintf("    app_name: %q", s.Name),
		"  logging:",
		`    driver: "json-file"`,
		"    options:",
		`      max-size: "10m"`,
		`      max-file: "3"`,
		`  network_mode: "host"`,
		"  environment:",
		fmt.Sprintf("    PORT: %d", s.Port),
	}
	lines = append(lines, s.envFileLines()...)
	lines = append(lines, s.healthcheckLines()...)
	lines = append(lines, "  restart: always")
	return indentLines(lines, unit)
}

func (s Service) envFileLines() []string {
	if s.EnvFile == "" {
		return nil
	}
	return []string{
		"  env_file:",
		"    - path: " + s.EnvFile,
		"      required: false",
	}
}

func (s Service) healthcheckLines() []string {
	if s.HealthPath == "" {
		return nil
	}
	url := fmt.Sprintf("http://localhost:%d%s", s.Port, s.HealthPath)
	return []string{
		"  healthcheck:",
		fmt.Sprintf(`    test: [ "CMD-SHELL", "wget -q --spider %s || exit 1" ]`, url),
		"    interval: 60s",
		"    timeout: 3s",
		"    retries: 10",
	}
}

// indentLines turns the 2-space levels of lines into levels of unit
// spaces, the service key itself being one level deep.
func indentLines(lines []string, unit int) string {
	var b strings.Builder
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " ")
		level := (len(line)-len(trimmed))/2 + 1
		b.WriteString(strings.Repeat(" ", unit*level) + trimmed + "\n")
	}
	return b.String()
}
