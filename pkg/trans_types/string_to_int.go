package trans_types

import (
	"strconv"
)

func AtoiOr0(nbr string) (int, error) {
	if nbr == "" {
		return 0, nil
	}
	it, err := strconv.Atoi(nbr)
	if err != nil {
		return 0, err
	}
	return it, err
}
