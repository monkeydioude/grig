package payload

import (
	customErr "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/pkg/errors"
)

func VerifyAndSanitizeCapybara(
	cp *model.Capybara,
) error {
	if cp == nil {
		return errors.Wrap(customErr.ErrNilPointer, "VerifyAndSanitizeCapybara: *model.Capybara")
	}
	if err := cp.Proxy.Verify(); err != nil {
		return errors.Wrap(err, "VerifyAndSanitizeCapybara: cp.Proxy.Verify")
	}
	cp.Sanitize()

	return nil
}

func VerifyAndSanitizeJosuke(jk *model.Josuke) error {
	if jk == nil {
		return errors.Wrap(customErr.ErrNilPointer, "VerifyAndSanitizeJosuke: *model.Josuke")
	}
	if err := jk.VerifyAndSanitize(); err != nil {
		return errors.Wrap(err, "VerifyAndSanitizeJosuke: *model.Josuke.VerifyAndSanitize")
	}
	return nil
}
