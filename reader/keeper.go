package reader

import (
	"sync"

	db "github.com/decentrio/ledger-reading/database/handlers"
	"github.com/decentrio/ledger-reading/lib/service"
	"github.com/stellar/go/support/log"
)

type Reader struct {
	service.BaseService

	db *db.DBHandler
	wg sync.WaitGroup
}
type ReaderOption func(*Reader)

func NewReader(
	// cfg *config.ExporterConfig,
	baseLogger *log.Entry,
	options ...ReaderOption,
) *Reader {
	r := &Reader{}
	r.BaseService = *service.NewBaseService("reader", r)
	for _, opt := range options {
		opt(r)
	}

	logger := baseLogger.WithField("module", "reader")
	logger.SetLevel(log.DebugLevel)
	r.BaseService.SetLogger(logger)

	r.db = db.NewDBHandler()

	return r
}

func (r *Reader) OnStart() error {
	r.Logger.Info("start services")
	r.wg.Add(1)

	go r.ledgerProcessing()
	return nil
}

func (r *Reader) OnStop() error {
	r.Logger.Info("stop services")
	r.wg.Wait()

	return nil
}
