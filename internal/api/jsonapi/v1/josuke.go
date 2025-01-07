package v1

import (
	"fmt"
	"monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/file"
	"net/http"
)

func (h Handler) JosukeSave(w http.ResponseWriter, r *http.Request, cp *model.Josuke) error {
	if r == nil || cp == nil {
		return fmt.Errorf("api.JosukeSave: %w", errors.ErrNilPointer)
	}
	cp.Path = h.Layout.ServerConfig.JosukeConfigPath
	cp.FileWriter = file.CreateAndWriteFile
	// Mutex is locked. Callback unlocks it
	defer (h.Layout.LockMutex())()
	if err := cp.Save(); err != nil {
		return fmt.Errorf("api.JosukeSave(): %w", err)
	}
	// replace layout's config
	// h.Layout.CapybaraConfig = &cp
	return nil
}
