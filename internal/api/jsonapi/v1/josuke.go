package v1

import (
	customErrors "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/internal/service/file"
	"monkeydioude/grig/internal/service/payload"
	"monkeydioude/grig/pkg/errors"
	"monkeydioude/grig/pkg/server/http_errors"
	"net/http"
)

func (h Handler) JosukeSave(w http.ResponseWriter, r *http.Request, jk *model.Josuke) error {
	if r == nil || jk == nil {
		return errors.Wrap(customErrors.ErrNilPointer, "api.JosukeSave")
	}
	jk.Path = h.Layout.ServerConfig.JosukeConfigPath
	jk.FileWriter = file.CreateAndWriteFile
	if err := payload.VerifyAndSanitizeJosuke(jk); err != nil {
		return http_errors.BadRequest(errors.Wrap(err, "api.JosukeSave"))
	}
	// Mutex is locked. Callback unlocks it
	defer (h.Layout.LockMutex())()
	if err := jk.Save(); err != nil {
		return errors.Wrap(err, "api.JosukeSave")
	}
	// replace layout's config
	// h.Layout.CapybaraConfig = &cp
	return nil
}
