package services

import (
	"fmt"
	"os"
)

func removeSymlink(dir, name string) error {
	return os.Remove(fmt.Sprintf("%s/%s", dir, name))
}
