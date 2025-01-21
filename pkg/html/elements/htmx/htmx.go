package htmx

import "monkeydioude/grig/internal/consts"

type Swap string

const (
	BeforeEnd   Swap = "beforeend"
	BeforeBegin Swap = "beforebegin"
	InnerHTML   Swap = "innerHTML"
)

func (s Swap) String() string {
	return string(s)
}

type Method string

const (
	Post Method = "post"
)

type Factory struct {
	Swap         Swap
	Target       string
	Ext          string
	SuccessMsg   string
	ErrorMsg     string
	IndexCounter string
	Method
}

func (f Factory) WithMessages(success, err string) Factory {
	f.SuccessMsg = success
	f.ErrorMsg = err
	return f
}

func (f Factory) WithSwapTarget(swap Swap, target string) Factory {
	f.Swap = swap
	f.Target = target
	return f
}

func (f Factory) WithIndexCounter(indexCounter string) Factory {
	f.IndexCounter = indexCounter
	return f
}

func NewFactory() Factory {
	return Factory{
		Swap:         InnerHTML,
		Target:       "this",
		SuccessMsg:   consts.FORM_SUCCESS_MSG,
		ErrorMsg:     consts.FORM_ERR_MSG,
		IndexCounter: "service-block",
	}
}

func NewJsonFactory() Factory {
	f := NewFactory()
	f.Ext = "json-enc-custom"
	return f
}
