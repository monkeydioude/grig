package page_data

type Page interface {
	Title() string
}

type Error string

func (e Error) Title() string {
	return string(e)
}
