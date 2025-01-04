package fs

import (
	"errors"
	"fmt"
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
	for _, file := range d.Files {
		if err := file.Save(); err != nil {
			return fmt.Errorf("fs.Dir.Save(): %w", err)
		}
	}

	return nil
}

func NewDirFromPathAndFileParser[F File](
	path string,
	parser func(path string) (F, error),
) (Dir[F], error) {
	dir := Dir[F]{
		Path:  path,
		mutex: &sync.Mutex{},
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return dir, fmt.Errorf("fs.NewDirFromPath: %w: %w", ErrReadDir, err)
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
