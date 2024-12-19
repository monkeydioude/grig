package fs

import (
	"os"
)

type File interface {
	Save() error
	Source() *os.File
}
