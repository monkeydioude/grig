package v1

import (
	"fmt"
	"log/slog"
	"monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/file"
	"monkeydioude/grig/internal/service/payload"
	"monkeydioude/grig/pkg/server/http_errors"
	"net/http"
	"strconv"
)

func (h Handler) CapybaraSave(w http.ResponseWriter, r *http.Request, logger *slog.Logger, cp *model.Capybara) error {
	if r == nil || cp == nil {
		return fmt.Errorf("api.CapybaraSave: %w", errors.ErrNilPointer)
	}

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return http_errors.BadRequest(fmt.Errorf("api.CapybaraSave(): invalid id: %w", err))
	}
	ref, ok := h.Layout.ServerConfig.CapybaraRefByID(id)
	if !ok {
		return http_errors.BadRequest(fmt.Errorf("api.CapybaraSave(): ID out of bounds"))
	}

	cp.Path = ref.Path
	cp.FileWriter = file.CreateAndWriteFile
	// the form posts only the fields grig displays, so carry over
	// everything else the file already holds
	if existing, err := file.UnmarshalFromPath[model.Capybara](ref.Path); err != nil {
		logger.Warn("api.CapybaraSave(): could not read existing config, keys it holds may be lost",
			"error", err, "path", ref.Path)
	} else {
		cp.MergeExtrasFrom(existing)
	}
	if err := payload.VerifyAndSanitizeCapybara(cp); err != nil {
		return http_errors.BadRequest(fmt.Errorf("api.CapybaraSave(): %w", err))
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
