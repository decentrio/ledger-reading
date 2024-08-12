package uploader

import (
	"sync"

	"github.com/decentrio/ledger-reading/lib/service"
	"github.com/stellar/go/support/log"
	"github.com/stellar/go/xdr"
)

type Uploader struct {
	service.BaseService

	// networkPassPhrase
	networkPassPhrase string
	// ledgerQueue channel for trigger new ledger
	ledgerReadChan <-chan xdr.LedgerCloseMeta

	TickerList            map[TokenPair]ITicker
	TickerListWithPoolKey map[string]ITicker
	TokenList             map[string]Token

	wg sync.WaitGroup
}

// ExporterOption sets an optional parameter on the State.
type UploaderOption func(*Uploader)

func NewUploader(
	// cfg *config.ExporterConfig,
	baseLogger *log.Entry,
	readChan <-chan xdr.LedgerCloseMeta,
	networkPassPhrase string,
	options ...UploaderOption,
) *Uploader {
	u := &Uploader{
		ledgerReadChan:        readChan,
		TickerListWithPoolKey: make(map[string]ITicker),
		TickerList:            make(map[TokenPair]ITicker),
		TokenList:             make(map[string]Token),
		networkPassPhrase:     networkPassPhrase,
	}

	u.BaseService = *service.NewBaseService("uploader", u)
	for _, opt := range options {
		opt(u)
	}

	logger := baseLogger.WithField("module", "uploader")
	logger.SetLevel(log.DebugLevel)
	u.BaseService.SetLogger(logger)

	return u
}

func (u *Uploader) OnStart() error {
	u.Logger.Info("start services")
	u.wg.Add(1)

	go u.ledgerProcessing()
	return nil
}

func (u *Uploader) OnStop() error {
	u.Logger.Info("stop services")
	u.wg.Wait()

	return nil

}