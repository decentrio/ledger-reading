package uploader

import (
	"fmt"
	"io"
	"time"

	"github.com/stellar/go/ingest"
	"github.com/stellar/go/xdr"
)

type Ledger struct {
	Hash         string `json:"hash,omitempty"`
	PrevHash     string `json:"prev_hash,omitempty"`
	Seq          uint32 `json:"seq,omitempty"`
	Transactions uint32 `json:"transaction,omitempty"`
	Operations   uint32 `json:"operations,omitempty"`
	LedgerTime   uint64 `json:"ledger_time,omitempty"`
}

// aggregation process
func (u *Uploader) ledgerProcessing() {
	defer u.wg.Done()
	for {
		select {
		// Terminate process
		case <-u.BaseService.Terminate():
			return
		// Receive a new tx
		case ledger := <-u.ledgerReadChan:
			u.handleReceiveNewLedger(ledger)
		default:
		}
		time.Sleep(time.Millisecond)
	}
}

// handle receive new ledger from exporter
func (u *Uploader) handleReceiveNewLedger(l xdr.LedgerCloseMeta) {
	ledger := getLedgerFromCloseMeta(l)
	// get tx
	fmt.Println(u.networkPassPhrase)
	txReader, err := ingest.NewLedgerTransactionReaderFromLedgerCloseMeta(u.networkPassPhrase, l)
	if err != nil {
		panic(err)
	}
	defer txReader.Close()

	var txs []*TransactionExtractor
	for {
		tx, err := txReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			u.Logger.Error(fmt.Sprintf("error txReader %s", err.Error()))
			continue
		}

		txWrapper := NewTransactionExtractor(tx, l.LedgerSequence(), ledger.LedgerTime)
		_, isIvkFn := txWrapper.IsInvokeHostFunctionTx()
		if txWrapper != nil && isIvkFn {
			txs = append(txs, txWrapper)
		}
	}

	for _, txExt:= range txs {
		LogData(txExt.Tx.Result.TransactionHash.HexString())
	}
}

func (u *Uploader) GetPoolLiquidity(tx *TransactionExtractor, ticker ITicker) (share, base, target, totalInUsd uint64) {
	contractData := tx.GetContractDataEntry()
	for _, cd := range contractData {
		if cd.ContractId == ticker.PoolContract {
			shareLiquidity, baseLiquidity, targetLiquidity, err := tx.GetPoolLiquidity(cd)
			if err != nil {
				u.Logger.Error(fmt.Sprintf("error get pool liquidity %s", err.Error()))
				break
			}
			share = shareLiquidity
			base = baseLiquidity
			target = targetLiquidity
		}
	}
	// calculate in usd
	baseCurrencyPriceInUsd := 1.0
	baseLiquidityInUsd := uint64(float64(base) * baseCurrencyPriceInUsd)

	targetLiquidityInUsd := uint64(float64(target))

	if baseLiquidityInUsd != 0 && targetLiquidityInUsd != 0 {
		totalInUsd = baseLiquidityInUsd + targetLiquidityInUsd
	} else if baseLiquidityInUsd == 0 && targetLiquidityInUsd != 0 {
		totalInUsd = 2 * targetLiquidityInUsd
	} else if baseLiquidityInUsd != 0 && targetLiquidityInUsd == 0 {
		totalInUsd = 2 * baseLiquidityInUsd
	}

	return share, base, target, totalInUsd
}

func getBaseTokenPrice(tradeType string, offerAmount uint64, returnAmount uint64) float64 {
	var price float64
	switch tradeType {
	case Buy:
		price = float64(offerAmount) / float64(returnAmount)
	case Sell:
		price = float64(returnAmount) / float64(offerAmount)
	}
	return price
}

func (u *Uploader) getTickerByTokenPair(offerTokenSorobanContract string, askedTokenSorobanContract string) (ITicker, string, bool) {
	// get token from soroban contract
	return ITicker{}, "", false
}

func getLedgerFromCloseMeta(ledgerCloseMeta xdr.LedgerCloseMeta) Ledger {
	var ledgerHeader xdr.LedgerHeaderHistoryEntry
	switch ledgerCloseMeta.V {
	case 0:
		ledgerHeader = ledgerCloseMeta.MustV0().LedgerHeader
	case 1:
		ledgerHeader = ledgerCloseMeta.MustV1().LedgerHeader
	default:
		panic(fmt.Sprintf("Unsupported LedgerCloseMeta.V: %d", ledgerCloseMeta.V))
	}

	timestamp := uint64(ledgerHeader.Header.ScpValue.CloseTime)

	return Ledger{
		Hash:       ledgerCloseMeta.LedgerHash().HexString(),
		PrevHash:   ledgerCloseMeta.PreviousLedgerHash().HexString(),
		Seq:        ledgerCloseMeta.LedgerSequence(),
		LedgerTime: timestamp,
	}
}
