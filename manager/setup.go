package manager

import (
	"github.com/stellar/go/support/log"

	"github.com/decentrio/ledger-reading/config"
)

func DefaultNewManager(cfg *config.ManagerConfig) *Manager {
	logger := log.New()
	return NewManager(cfg, logger)
}
