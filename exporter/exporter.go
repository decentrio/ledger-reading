package exporter

import (
	"context"
	"sync"
	"time"

	backends "github.com/stellar/go/ingest/ledgerbackend"
	"github.com/stellar/go/support/log"
	"github.com/stellar/go/xdr"

	"github.com/decentrio/ledger-reading/config"
	"github.com/decentrio/ledger-reading/lib/service"
)

const (
	QueueSize          = 10000
	DefaultPrepareStep = 64
)

const (
	TickerListUrl = "https://raw.githubusercontent.com/decentrio/token-list/main/ticker_lists.json"
)

type Exporter struct {
	service.BaseService

	Config *config.ExporterConfig

	ctx               context.Context
	backend           backends.LedgerBackend
	networkPassPhrase string

	// ledgerQueue channel for trigger new ledger
	ledgerQueue chan xdr.LedgerCloseMeta

	// isSync is flag represent if services is
	// re-synchronize
	isSync      bool
	prepareStep uint32

	StartLedgerSeq uint32
	CurrLedgerSeq  uint32

	wg sync.WaitGroup
}

// ExporterOption sets an optional parameter on the State.
type ExporterOption func(*Exporter)

func NewExporter(
	cfg *config.ExporterConfig,
	baseLogger *log.Entry,
	options ...ExporterOption,
) *Exporter {
	e := &Exporter{
		ledgerQueue: make(chan xdr.LedgerCloseMeta, QueueSize),
		prepareStep: DefaultPrepareStep,
		isSync:      false,
		Config:      cfg,
	}

	e.BaseService = *service.NewBaseService("Exporter", e)
	for _, opt := range options {
		opt(e)
	}

	logger := baseLogger.WithField("module", "exporter")
	logger.SetLevel(log.ErrorLevel)
	e.BaseService.SetLogger(logger)

	e.StartLedgerSeq = e.Config.StartLedgerHeight
	e.CurrLedgerSeq = e.Config.CurrLedgerHeight

	e.ctx = context.Background()
	e.backend, e.networkPassPhrase = newLedgerBackend(e.ctx, *e.Config, e.Logger)

	return e
}

func (e *Exporter) OnStart() error {
	e.Logger.Info("start services")
	e.wg.Add(1)
	// Note that when using goroutines, you need to be careful to ensure that no
	// race conditions occur when accessing the txQueue.
	go e.aggregation()
	return nil
}

func (e *Exporter) OnStop() error {
	e.Logger.Info("stop services")
	e.backend.Close()
	e.wg.Wait()

	return nil
}

func (e *Exporter) aggregation() {
	defer e.wg.Done()
	for {
		select {
		// Terminate process
		case <-e.BaseService.Terminate():
			return
		default:
			e.getNewLedger()
		}
		time.Sleep(time.Millisecond)
	}
}

func (e *Exporter) GetLedgerChanPipe() <-chan xdr.LedgerCloseMeta {
	return e.ledgerQueue
}

func (e *Exporter) GetNetworkPassPhrase() string {
	return e.networkPassPhrase
}
