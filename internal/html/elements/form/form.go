package form

import "fmt"

type Type string

const (
	String Type = "string"
	Number Type = "number"
)

func (t Type) String() string {
	return string(t)
}

type FormGroup interface {
	ID() string
	Placeholder() string
	Value() string
	Type() Type
	Required() bool
}

type BaseStringFormGroup struct {
	id          string
	placeholder string
	required    bool
}

type stringFormGroup struct {
	BaseStringFormGroup
	value string
}

func StringFormGroup(id, placeholder, value string) stringFormGroup {
	return stringFormGroup{
		BaseStringFormGroup: BaseStringFormGroup{
			id:          id,
			placeholder: placeholder,
		},
		value: value,
	}
}

func (s stringFormGroup) ID() string {
	return s.id
}

func (s stringFormGroup) Placeholder() string {
	return s.placeholder
}

func (s stringFormGroup) Value() string {
	return s.value
}

func (s stringFormGroup) Type() Type {
	return String
}

func (s stringFormGroup) Required() bool {
	return s.required
}

func (s stringFormGroup) NotRequired() stringFormGroup {
	s.required = false
	return s
}

type numberFormGroup struct {
	BaseStringFormGroup
	value fmt.Stringer
}

type EmptyStringer struct{}

func (EmptyStringer) String() string {
	return ""
}

func NumberFormGroup(id, placeholder string, value fmt.Stringer) numberFormGroup {
	if value == nil {
		value = EmptyStringer{}
	}
	return numberFormGroup{
		BaseStringFormGroup: BaseStringFormGroup{
			id:          id,
			placeholder: placeholder,
			required:    true,
		},
		value: value,
	}
}

func (n numberFormGroup) NotRequired() numberFormGroup {
	n.required = false
	return n
}

func (n numberFormGroup) ID() string {
	return n.id
}

func (n numberFormGroup) Placeholder() string {
	return n.placeholder
}

func (n numberFormGroup) Value() string {
	return n.value.String()
}

func (n numberFormGroup) Type() Type {
	return Number
}

func (n numberFormGroup) Required() bool {
	return n.required
}
