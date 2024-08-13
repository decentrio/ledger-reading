package manager

import (
	"github.com/stellar/go/support/log"
)

func DefaultNewManager(from, to int32) *Manager {
	logger := log.New()
	return NewManager(logger, from, to)
}
