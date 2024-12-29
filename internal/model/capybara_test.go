package model

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestICanSaveACapybaraFile(t *testing.T) {
	cap := Capybara{
		Proxy: Proxy{
			Port:    12,
			TLSHost: "cabane123.com",
		},
		Services: []ServiceDefinition{
			{
				ID:       "1",
				Method:   "test_1",
				Pattern:  "string",
				Port:     2222,
				Protocol: "http",
			},
		},
		Path: "/path/test",
	}
	goal, err := json.Marshal(cap)
	assert.NoError(t, err)
	fw := func(path string, data []byte, mode os.FileMode) error {
		assert.Equal(t, "/path/test", path)
		assert.Equal(t, goal, data)

		um := Capybara{}
		assert.NoError(t, json.Unmarshal(data, &um))
		assert.Nil(t, um.FileWriter)
		assert.Empty(t, um.Path)
		return nil
	}
	cap.FileWriter = fw

}
