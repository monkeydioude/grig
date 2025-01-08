package payload

import (
	"log"
	customErr "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/model"
	"monkeydioude/grig/pkg/errors"
	"slices"
)

func VerifyAndSanitizeCapybara(
	cp *model.Capybara,
) error {
	if cp == nil {
		return errors.Wrap(customErr.ErrNilPointer, "VerifyAndSanitizeCapybara: *model.Capybara")
	}
	if err := cp.Proxy.Verify(); err != nil {
		return errors.Wrap(err, "VerifyAndSanitizeCapybara: *model.Capybara")
	}
	// use of a custom index, so we dont range into a non existant element
	// in case we delete one
	i := 0
	for i < len(cp.Services) {
		sd := cp.Services[i]
		if err := sd.Verify(); err != nil {
			// no need to throw an error here, we just remove the element
			log.Printf("[ERR ] VerifyAndSanitizeCapybara: %+v", err.Error())
			cp.Services = slices.Delete(cp.Services, i, i+1)
		} else {
			i++
		}
	}

	return nil
}

// func VerifyAndSanitizeJosuke(jk *model.Josuke) error {
// 	if jk == nil {
// 		return pkgErr.Wrap(errors.ErrNilPointer, "VerifyAndSanitizeJosuke: *model.Capybara")
// 	}
// }
