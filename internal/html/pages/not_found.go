package pages

type Error string

func (e Error) Title() string {
	return string(e)
}
