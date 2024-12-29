package v1

import (
	"fmt"
	"io"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/action"
	"monkeydioude/grig/internal/service/fs"
	"net/http"
)

func (h Handler) CapybaraSave(w http.ResponseWriter, r *http.Request) error {
	var data []byte
	{
		var err error
		data, err = io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("api.CapybaraCreate(): %w", err)
		}
	}
	// create a base Capybara config
	// cp := h.Layout.CapybaraConfig.CloneBase()
	// hydrate this Capybara entity by populating it with unmarshaled data
	cp := model.Capybara{
		Path:       h.Layout.ServerConfig.CapybaraConfigPath,
		FileWriter: fs.CreateAndWriteFile,
	}
	if err := action.HydrateCapybaraFromPayload(data, &cp); err != nil {
		return fmt.Errorf("api.CapybaraCreate(): %w", err)
	}
	// Mutex is locked. Callback unlocks it
	defer (h.Layout.LockMutex())()
	if err := cp.Save(); err != nil {
		return fmt.Errorf("api.CapybaraCreate(): %w", err)
	}
	// replace layout's config
	// h.Layout.CapybaraConfig = &cp
	return nil
}
