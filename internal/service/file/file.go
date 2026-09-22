package file

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"monkeydioude/grig/internal/errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type File interface {
	Save() error
}

func AppendToThisFileDirectory(appendThisFilesPath, toThisFileDir string) string {
	dir := filepath.Dir(toThisFileDir)
	appendedPath := filepath.Join(dir, appendThisFilesPath)
	res, err := filepath.Abs(appendedPath)
	if err != nil {
		slog.Warn("AppendToThisFileDirectory: filepath.Abs", "error", err)
		return appendThisFilesPath
	}
	return res
}

// codec is the (un)marshaling pair used for a config file format.
type codec struct {
	unmarshal func([]byte, any) error
	marshal   func(any) ([]byte, error)
}

// yamlMarshalIndent encodes to YAML with 2-space indentation, which is
// the conventional style. yaml.Marshal defaults to 4 spaces.
func yamlMarshalIndent(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// codecForPath picks the codec matching the file extension,
// defaulting to JSON for any unknown extension.
func codecForPath(path string) codec {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".yaml", ".yml":
		return codec{unmarshal: yaml.Unmarshal, marshal: yamlMarshalIndent}
	default:
		return codec{unmarshal: json.Unmarshal, marshal: json.Marshal}
	}
}

// MarshalForPath encodes data in the format matching path's extension.
func MarshalForPath(path string, data any) ([]byte, error) {
	raw, err := codecForPath(path).marshal(data)
	if err != nil {
		return nil, fmt.Errorf("fs.MarshalForPath: %q: %s: %s", path, errors.ErrMarshaling, err)
	}
	return raw, nil
}

func UnmarshalFromPath[F File](configPath string) (F, error) {
	rawData, err := os.ReadFile(configPath)
	var res F
	if err != nil {
		return res, fmt.Errorf("fs.UnmarshalFromPath: %q: %s: %s", configPath, errors.ErrReadingFile, err)
	}
	err = codecForPath(configPath).unmarshal(rawData, &res)
	if err != nil {
		return res, fmt.Errorf("fs.UnmarshalFromPath: %q: %s: %s", configPath, errors.ErrUnmarshaling, err)
	}
	return res, nil
}

func CreateAndWriteFile(path string, data []byte, mode os.FileMode) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("fs.CreateAndWriteFile(): %q: %w:, %w", path, errors.ErrCreatingFile, err)
		}
		file.Close()
	} else if err != nil {
		return fmt.Errorf("fs.CreateAndWriteFile(): %q: %w: %w", path, errors.ErrCheckingFile, err)
	}
	if err := os.WriteFile(path, data, mode); err != nil {
		return fmt.Errorf("fs.CreateAndWriteFile(): %q: %w: %w", path, errors.ErrWritingFile, err)
	}
	return nil
}
