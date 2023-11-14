package errwrap

import (
	"errors"
	"fmt"
)

type ErrWrap struct {
	err  error
	werr error
}

func NewErrWrap(err error, werr error) error {
	return &ErrWrap{err, werr}
}

func (ew *ErrWrap) Error() string {
	return fmt.Sprintf("%s: %s", ew.err, ew.werr)
}

func (ew *ErrWrap) Unwrap() error {
	return ew.werr
}

func (ew *ErrWrap) Is(target error) bool {
	return errors.Is(ew.err, target) ||
		errors.Is(ew.werr, target)
}
