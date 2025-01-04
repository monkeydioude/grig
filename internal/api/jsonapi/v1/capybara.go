package v1

import (
	"fmt"
	"monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/action"
	"monkeydioude/grig/pkg/fs"
	"net/http"
)

func (h Handler) CapybaraSave(w http.ResponseWriter, r *http.Request, cp *model.Capybara) error {
	if r == nil || cp == nil {
		return fmt.Errorf("api.CapybaraSave: %w", errors.ErrNilPointer)
	}
	cp.Path = h.Layout.ServerConfig.CapybaraConfigPath
	cp.FileWriter = fs.CreateAndWriteFile
	if err := action.HydrateCapybaraFromPayload(cp); err != nil {
		return errors.BadRequest(fmt.Errorf("api.CapybaraSave(): %w", err))
	}
	// Mutex is locked. Callback unlocks it
	defer (h.Layout.LockMutex())()
	if err := cp.Save(); err != nil {
		return fmt.Errorf("api.CapybaraSave(): %w", err)
	}
	// replace layout's config
	// h.Layout.CapybaraConfig = &cp
	return nil
}
