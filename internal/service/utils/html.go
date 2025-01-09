package utils

import (
	"fmt"
	"monkeydioude/grig/pkg/model"
)

func GetInputName(hp model.IndexBuilder, typ string) string {
	if hp == nil {
		return ""
	}
	res := ""
	for hp.GetParent() != nil {
		res = fmt.Sprintf("[%s][%d]%s", hp.GetName(), hp.GetIndex(), res)
		hp = hp.GetParent()
	}
	if typ != "" {
		typ = fmt.Sprintf("[%s]", typ)
	}
	return fmt.Sprintf("%s[%d]%s%s", hp.GetName(), hp.GetIndex(), res, typ)
}
