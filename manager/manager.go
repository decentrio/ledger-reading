package manager

import (
	"time"

	"github.com/decentrio/ledger-reading/lib/service"
	lreader "github.com/decentrio/ledger-reading/reader"
	"github.com/stellar/go/support/log"
)

// Manager is the root service that manage all services
type Manager struct {
	service.BaseService

	// sub services that are controlled by manager services
	r *lreader.Reader
}

const (
	PaddingLedger = 2560
)

// StateOption sets an optional parameter on the State.
type ManagerOption func(*Manager)

// NewBaseService creates a new manager.
func NewManager(
	baseLogger *log.Entry,
	options ...ManagerOption,
) *Manager {
	m := &Manager{}

	// initialize exporter sub services
	m.r = lreader.NewReader(baseLogger)

	m.BaseService = *service.NewBaseService("manager", m)
	for _, opt := range options {
		opt(m)
	}

	m.BaseService.SetLogger(baseLogger.WithField("module", "manager"))

	return m
}

func (m *Manager) OnStart() error {
	m.Logger.Info("start services")

	// start uploader services
	if err := m.r.Start(); err != nil {
		return err
	}

	return nil
}

func (m *Manager) OnStop() error {
	m.Logger.Info("stop services")
	m.r.Stop()

	time.Sleep(time.Second)
	return nil
}
