package model

import "os"

type Capybara struct {
}

func (c Capybara) Save() error {
	return nil
}

func (c Capybara) Source() *os.File {
	return nil
}
