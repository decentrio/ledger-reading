package importer

import (
	"sync"

	"github.com/decentrio/ledger-reading/lib/service"
	"github.com/stellar/go/support/log"
)

type Importer struct {
	service.BaseService

	updateTime uint8

	tokenListUpdateCb  func(token Token)
	tickerListUpdateCb func(ticker Ticker)

	TickerList map[string]Ticker
	TokenList  map[string]Token

	wg sync.WaitGroup
}

// ExporterOption sets an optional parameter on the State.
type ImporterOption func(*Importer)

func NewImporter(
	baseLogger *log.Entry,
	tokenCb func(token Token),
	tickerCb func(ticker Ticker),
	options ...ImporterOption,
) *Importer {
	i := &Importer{
		tokenListUpdateCb:  tokenCb,
		tickerListUpdateCb: tickerCb,
	}

	i.BaseService = *service.NewBaseService("importer", i)
	for _, opt := range options {
		opt(i)
	}

	logger := baseLogger.WithField("module", "importer")
	logger.SetLevel(log.DebugLevel)
	i.BaseService.SetLogger(logger)

	return i
}

func (i *Importer) OnStart() error {
	i.Logger.Info("start services")
	i.wg.Add(2)

	go i.FetchTickerList()
	go i.FetchTokenList()

	return nil
}

func (i *Importer) OnStop() error {
	i.Logger.Info("stop services")
	i.wg.Wait()

	return nil
}
