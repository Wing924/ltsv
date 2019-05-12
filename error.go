package ltsv

import (
	"errors"
)

var (
	ErrMissingLabel = errors.New("missing label")
	ErrInvalidLabel = errors.New("invalid label")
)
