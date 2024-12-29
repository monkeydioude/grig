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
		return fmt.Errorf("AssertCapybaraPayload: *model.Capybara: %w", errors.ErrNilPointer)
	}
	if err := json.Unmarshal(data, cp); err != nil {
		return fmt.Errorf("AssertCapybaraPayload: %w", err)
	}
	if err := cp.Verify(); err != nil {
		return fmt.Errorf("AssertCapybaraPayload: %w", err)
	}
	return nil
}
