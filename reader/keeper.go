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
	FromSeq int32
	ToSeq  int32
	wg sync.WaitGroup
}
type ReaderOption func(*Reader)

func NewReader(
	// cfg *config.ExporterConfig,
	baseLogger *log.Entry,
	from int32,
	to int32,
	options ...ReaderOption,
) *Reader {
	r := &Reader{
		FromSeq: from,
		ToSeq: to,
	}
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

	go r.LedgerProcessing()
	return nil
}

func (r *Reader) OnStop() error {
	r.Logger.Info("stop services")
	r.wg.Done()
	r.wg.Wait()

	return nil
}
