package manager

import (
	"github.com/stellar/go/support/log"
)

func DefaultNewManager() *Manager {
	logger := log.New()
	return NewManager(logger, 0, 0)
}

func NewManagerFromTo(from, to int32) *Manager {
	logger := log.New()
	return NewManager(logger, from, to)
}
