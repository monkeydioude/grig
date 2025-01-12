package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandVerifyAndSanitize(t *testing.T) {
	cmds := Command{
		Parts:   []string{"a", "", "b"},
		parent:  nil,
		Indexer: Indexer{},
	}

	assert.NoError(t, cmds.VerifyAndSanitize())
	assert.Equal(t, []string{"a", "b"}, cmds.Parts)
}
