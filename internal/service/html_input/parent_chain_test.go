package htmlinput

import (
	"monkeydioude/grig/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFillIndexesFromHTMLInput(t *testing.T) {
	trial := "deployment[1][branch][2][action][3]"

	cmd := model.Command{}
	cmd.InitParent()
	assert.NoError(t, FillIndexesFromHTMLInput(cmd.GetParent(), trial))
	act := cmd.GetParent()
	assert.Equal(t, 3, act.GetIndex())
	br := act.GetParent()
	assert.Equal(t, 2, br.GetIndex())
	dep := br.GetParent()
	assert.Equal(t, 1, dep.GetIndex())
}
