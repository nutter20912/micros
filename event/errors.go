package event

import (
	"errors"
)

type EventError error

var (
	ErrMessageExpired    = errors.New("message is expired")
	ErrMessageConflicted = errors.New("message is conflicted")
)
