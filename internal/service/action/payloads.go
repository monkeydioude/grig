package action

import (
	"encoding/json"
	"fmt"
	"monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/model"
)

func HydrateCapybaraFromPayload(
	data []byte,
	cp *model.Capybara,
) error {
	if cp == nil {
		return fmt.Errorf("HydrateCapybaraFromPayload: *model.Capybara: %w", errors.ErrNilPointer)
	}
	if err := json.Unmarshal(data, cp); err != nil {
		return fmt.Errorf("HydrateCapybaraFromPayload: %w", err)
	}
	if err := cp.Verify(); err != nil {
		return fmt.Errorf("HydrateCapybaraFromPayload: %w", err)
	}
	return nil
}
