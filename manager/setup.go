package manager

import (
	"github.com/stellar/go/support/log"
)

func DefaultNewManager() *Manager {
	logger := log.New()
	return NewManager(logger)
}
