package fs

import (
	"encoding/json"
	"fmt"
	"log"
	"monkeydioude/grig/internal/errors"
	"os"
	"path/filepath"
)

type File interface {
	Save() error
	Source() *os.File
}

func AppendToThisFileDirectory(appendThisFilesPath, toThisFileDir string) string {
	dir := filepath.Dir(toThisFileDir)
	appendedPath := filepath.Join(dir, appendThisFilesPath)
	res, err := filepath.Abs(appendedPath)
	if err != nil {
		log.Printf("[WARN] AppendToThisFileDirectory: filepath.Abs: %s\n", err)
		return appendThisFilesPath
	}
	return res
}

func UnmarshalFromPath[F File](josukeConfigPath string) (F, error) {
	rawData, err := os.ReadFile(josukeConfigPath)
	var res F
	if err != nil {
		return res, fmt.Errorf("josuke.UnmarshalFromPath: %s: %s", errors.ErrReadingFile, err)
	}
	err = json.Unmarshal(rawData, &res)
	if err != nil {
		return res, fmt.Errorf("josuke.UnmarshalFromPath: %s: %s", errors.ErrUnmarshaling, err)
	}
	return res, nil
}
