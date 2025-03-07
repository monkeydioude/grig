package file

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"monkeydioude/grig/internal/errors"
	"os"
	"path/filepath"
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

func UnmarshalFromPath[F File](configPath string) (F, error) {
	rawData, err := os.ReadFile(configPath)
	var res F
	if err != nil {
		return res, fmt.Errorf("fs.UnmarshalFromPath: %s: %s", errors.ErrReadingFile, err)
	}
	err = json.Unmarshal(rawData, &res)
	if err != nil {
		return res, fmt.Errorf("fs.UnmarshalFromPath: %s: %s", errors.ErrUnmarshaling, err)
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
