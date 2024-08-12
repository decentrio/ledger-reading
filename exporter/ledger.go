package exporter

import (
	"fmt"

	backends "github.com/stellar/go/ingest/ledgerbackend"
	"github.com/stellar/go/xdr"
)

func (e *Exporter) getNewLedger() {
	// prepare range
	from, to := e.prepare()
	// get ledger
	if !e.isSync {
		for seq := from; seq < to; seq++ {
			ledgerCloseMeta, err := e.backend.GetLedger(e.ctx, seq)
			if err != nil {
				e.Logger.Error(fmt.Sprintf("error get ledger %s", err.Error()))
				return
			}

			go func(l xdr.LedgerCloseMeta) {
				e.ledgerQueue <- l
			}(ledgerCloseMeta)
		}
	} else {
		seq := e.StartLedgerSeq
		ledgerCloseMeta, err := e.backend.GetLedger(e.ctx, seq)
		if err != nil {
			e.Logger.Error(fmt.Sprintf("error get ledger %s", err.Error()))
			return
		}

		go func(l xdr.LedgerCloseMeta) {
			e.ledgerQueue <- l
		}(ledgerCloseMeta)
		e.StartLedgerSeq++
	}
}

func (as *Exporter) prepare() (uint32, uint32) {
	if !as.isSync {
		from := as.StartLedgerSeq
		to := from + 1

		var ledgerRange backends.Range
		if to > as.CurrLedgerSeq {
			ledgerRange = backends.UnboundedRange(from)
		} else {
			ledgerRange = backends.BoundedRange(from, to)
		}

		fmt.Println(ledgerRange)
		err := as.backend.PrepareRange(as.ctx, ledgerRange)
		if err != nil {
			as.Logger.Errorf("error prepare %s", err.Error())
			return 0, 0 // if prepare error, we should skip here
		} else {
			if to > as.CurrLedgerSeq {
				as.isSync = true
			}
		}
		as.StartLedgerSeq += DefaultPrepareStep
		return from, to
	}

	return 0, 0
}
