package action

import (
	"fmt"
	"monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/model"
)

func HydrateCapybaraFromPayload(
	cp *model.Capybara,
) error {
	if cp == nil {
		return fmt.Errorf("HydrateCapybaraFromPayload: *model.Capybara: %w", errors.ErrNilPointer)
	}
	if err := cp.Verify(); err != nil {
		return fmt.Errorf("HydrateCapybaraFromPayload: %w", err)
	}
	return nil
}
