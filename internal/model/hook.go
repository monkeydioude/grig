package model

type Hook struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Path   string `json:"path"`
	Secret string `json:"secret"`
}

func (Hook) GetName() string {
	return "hook"
}
