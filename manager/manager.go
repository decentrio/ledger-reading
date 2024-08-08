package manager

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/decentrio/ledger-reading/config"
	"github.com/decentrio/ledger-reading/exporter"
	"github.com/decentrio/ledger-reading/importer"
	"github.com/decentrio/ledger-reading/lib/service"
	"github.com/decentrio/ledger-reading/uploader"
	"github.com/stellar/go/support/log"
)

// Manager is the root service that manage all services
type Manager struct {
	service.BaseService

	// config of Manager
	cfg *config.ManagerConfig

	// sub services that are controlled by manager services
	e *exporter.Exporter
	i *importer.Importer
	u *uploader.Uploader
}

const (
	PaddingLedger = 2560
)

// StateOption sets an optional parameter on the State.
type ManagerOption func(*Manager)

// NewBaseService creates a new manager.
func NewManager(
	cfg *config.ManagerConfig,
	baseLogger *log.Entry,
	options ...ManagerOption,
) *Manager {
	m := &Manager{
		cfg: cfg,
	}

	// initialize exporter sub services
	m.e = exporter.NewExporter(cfg.ExporterConfig, baseLogger)
	readChan := m.e.GetLedgerChanPipe()
	networkPassPhrase := m.e.GetNetworkPassPhrase()

	// initlialize uploader sub services
	m.u = uploader.NewUploader(baseLogger, readChan, networkPassPhrase)

	// initialize importer sub services
	newTokenCb := func(token importer.Token) { m.u.UploadNewToken(token) }
	newTickerCb := func(ticker importer.Ticker) { m.u.UploadNewTicker(ticker) }
	m.i = importer.NewImporter(baseLogger, newTokenCb, newTickerCb)

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
	if err := m.u.Start(); err != nil {
		return err
	}

	// start importer services
	if err := m.i.Start(); err != nil {
		return err
	}

	// start exporter services
	if err := m.e.Start(); err != nil {
		return err
	}

	return nil
}

func (m *Manager) OnStop() error {
	m.Logger.Info("stop services")
	m.e.Stop()

	// save current config
	asConfig := *m.e.Config
	asConfig.StartLedgerHeight = m.e.StartLedgerSeq - PaddingLedger

	bz, err := json.Marshal(asConfig)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(m.cfg.AggregationConfigFile())
	err = config.WriteState(m.cfg.AggregationConfigFile(), bz, 0777)
	if err != nil {
		fmt.Println(err.Error())
	}

	time.Sleep(time.Second)
	return nil
}
