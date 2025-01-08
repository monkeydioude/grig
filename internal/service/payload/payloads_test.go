package payload

import (
	customErr "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestICanVerifyAndSanitizeCapybara(t *testing.T) {
	cp := model.Capybara{
		Proxy: model.Proxy{Port: 80, TLSHost: "test.aye"},
		Services: []model.ServiceDefinition{
			{
				ID:      "1",
				Method:  "string",
				Pattern: "/lets/go",
				Port:    81,
			},
			{
				ID:      "",
				Method:  "",
				Pattern: "",
				Port:    0,
			},
		},
	}

	goal := model.Capybara{
		Proxy: model.Proxy{Port: 80, TLSHost: "test.aye"},
		Services: []model.ServiceDefinition{
			{
				ID:      "1",
				Method:  "string",
				Pattern: "/lets/go",
				Port:    81,
			},
		},
	}
	assert.NoError(t, VerifyAndSanitizeCapybara(&cp))
	assert.Equal(t, goal, cp)
}

func TestICanNOTVerifyAndSanitizeCapybaraOnWrongProxy(t *testing.T) {
	cp := model.Capybara{
		Proxy: model.Proxy{},
		Services: []model.ServiceDefinition{
			{
				ID:      "1",
				Method:  "string",
				Pattern: "/lets/go",
				Port:    81,
			},
			{
				ID:      "",
				Method:  "",
				Pattern: "",
				Port:    0,
			},
		},
	}
	assert.ErrorIs(t, VerifyAndSanitizeCapybara(&cp), customErr.ErrModelVerifyInvalidValue)
}

func TestICanNOTVerifyAndSanitizeCapybaraOnNilPtr(t *testing.T) {
	var cp *model.Capybara
	assert.ErrorIs(t, VerifyAndSanitizeCapybara(cp), customErr.ErrNilPointer)
}
