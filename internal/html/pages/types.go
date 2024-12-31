package pages

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
	Label() string
	Placeholder() string
	Value() string
	Type() Type
}

type BaseStringFormGroup struct {
	id          string
	label       string
	placeholder string
}

type StringFormGroup struct {
	BaseStringFormGroup
	value string
}

func NewStringFormGroup(id, label, placeholder, value string) StringFormGroup {
	return StringFormGroup{
		BaseStringFormGroup: BaseStringFormGroup{
			id:          id,
			label:       label,
			placeholder: placeholder,
		},
		value: value,
	}
}

func (s StringFormGroup) ID() string {
	return s.id
}

func (s StringFormGroup) Label() string {
	return s.label
}

func (s StringFormGroup) Placeholder() string {
	return s.placeholder
}

func (s StringFormGroup) Value() string {
	return s.value
}

func (s StringFormGroup) Type() Type {
	return String
}

type NumberFormGroup struct {
	BaseStringFormGroup
	value fmt.Stringer
}

func NewNumberFormGroup(id, label, placeholder string, value fmt.Stringer) NumberFormGroup {
	return NumberFormGroup{
		BaseStringFormGroup: BaseStringFormGroup{
			id:          id,
			label:       label,
			placeholder: placeholder,
		},
		value: value,
	}
}

func (n NumberFormGroup) ID() string {
	return n.id
}

func (n NumberFormGroup) Label() string {
	return n.label
}

func (n NumberFormGroup) Placeholder() string {
	return n.placeholder
}

func (n NumberFormGroup) Value() string {
	return n.value.String()
}

func (n NumberFormGroup) Type() Type {
	return Number
}
