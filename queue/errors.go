package queue

import (
	"errors"
)

type EventError error

var (
	ErrNoMetadataReceived = errors.New("no metadata received")
	ErrMessageExpired     = errors.New("message is expired")
	ErrMessageConflicted  = errors.New("message is conflicted")
)
