package utils

import (
	"fmt"
	"monkeydioude/grig/internal/model"
)

func GetInputName(hp model.IndexBuilder) string {
	if hp == nil {
		return ""
	}
	res := ""
	for hp.GetParent() != nil {
		res = fmt.Sprintf("[%s][%d]%s", hp.GetName(), hp.GetIndex(), res)
		hp = hp.GetParent()
	}
	return fmt.Sprintf("%s[%d]%s", hp.GetName(), hp.GetIndex(), res)
}
