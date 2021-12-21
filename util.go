package hotomata

import (
	"errors"
	"fmt"
)

func newError(msg string, params ...interface{}) error {
	return errors.New(fmt.Sprintf(msg, params...))
}
