package fs

import (
	"encoding/json"
	"errors"
	"fmt"
	customErrors "monkeydioude/grig/internal/errors"
	"os"
	"path/filepath"
	"sync"
)

type Dir[F File] struct {
	Path  string
	Files map[string]F
	mutex *sync.Mutex
}

func (d Dir[F]) Save() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	data, err := json.Marshal(d.Files)
	if err != nil {
		return fmt.Errorf("fs.Dir.Save(): %w: %w", customErrors.ErrMarshaling, err)
	}
	if err := os.WriteFile(d.Path, data, os.ModePerm); err != nil {
		return fmt.Errorf("fs.Dir.Save(): %w: %w", customErrors.ErrWritingFile, err)
	}
	return nil
}

func NewDirFromPathAndFileParser[F File](path string, parser func(path string) (F, error)) (Dir[F], error) {
	dir := Dir[F]{
		Path: path,
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return dir, fmt.Errorf("fs.NewDirFromPath: %w: %w", customErrors.ErrReadDir, err)
	}
	var stackedErrs error
	dir.Files = make(map[string]F, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		// content, err := os.ReadFile(filepath.Join(path, entry.Name()))
		// if err != nil {
		// 	stackedErrs = errors.Join(stackedErrs, fmt.Errorf("fs.NewDirFromPath: ReadFile: %w: %w", customErrors.ErrReadEntryFile, err))
		// 	continue
		// }

		file, err := parser(filepath.Join(path, entry.Name()))
		if err != nil {
			stackedErrs = errors.Join(stackedErrs, fmt.Errorf("fs.NewDirFromPath: NewServiceFromPath: %w", err))
			continue
		}
		dir.Files[entry.Name()] = file
	}
	return dir, stackedErrs
}
