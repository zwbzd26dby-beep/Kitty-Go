package types

import "fmt"

func errEmptyContent(kind string) error {
	return fmt.Errorf("%s content must not be empty", kind)
}
