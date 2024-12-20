package fs

import (
	"errors"
	"fmt"
	customErrors "monkeydioude/grig/internal/errors"
	"os"
	"path/filepath"
)

type Dir[F File] struct {
	Path  string
	Files map[string]F
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
