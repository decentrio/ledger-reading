package uploader

import (
	"github.com/stellar/go/ingest"
	"github.com/stellar/go/xdr"
)

const (
	SUCCESS = "success"
	FAILED  = "failed"
)

type IndexerTransactionExtractor struct {
	LedgerSequence uint32
	Tx             ingest.LedgerTransaction
	Ops            []transactionOperationWrapper
	Time           uint64
}

type ActionType string

const (
	Unkown           ActionType = "unknown"
	Swap             ActionType = "swap"
	ProvideLiquidity ActionType = "provide_liquidity"
)

func NewIndexerTransactionExtractor(tx ingest.LedgerTransaction, ledgerSeq uint32, processedUnixTime uint64) *IndexerTransactionExtractor {
	var ops []transactionOperationWrapper
	for opi, op := range tx.Envelope.Operations() {
		operation := transactionOperationWrapper{
			index:          uint32(opi),
			txIndex:        tx.Index,
			operation:      op,
			ledgerSequence: ledgerSeq,
		}

		ops = append(ops, operation)	
	}

	return &IndexerTransactionExtractor{
		LedgerSequence: ledgerSeq,
		Tx:             tx,
		Ops:            ops,
		Time:           processedUnixTime,
	}
}


func (tw *IndexerTransactionExtractor) GetTransactionHash() string {
	return tw.Tx.Result.TransactionHash.HexString()
}

func (tw *IndexerTransactionExtractor) GetStatus() string {
	if tw.Tx.Result.Successful() {
		return SUCCESS
	}

	return FAILED
}

func (tw *IndexerTransactionExtractor) GetLedgerSequence() uint32 {
	return tw.LedgerSequence
}

func (tw *IndexerTransactionExtractor) GetApplicationOrder() uint32 {
	return tw.Tx.Index
}

func (tw *IndexerTransactionExtractor) GetEnvelopeXdr() []byte {
	bz, _ := tw.Tx.Envelope.MarshalBinary()
	return bz
}

func (tw *IndexerTransactionExtractor) GetResultXdr() []byte {
	bz, _ := tw.Tx.Result.MarshalBinary()
	return bz
}

func (tw *IndexerTransactionExtractor) GetResultMetaXdr() []byte {
	txResultMeta := xdr.TransactionResultMeta{
		Result:            tw.Tx.Result,
		FeeProcessing:     tw.Tx.FeeChanges,
		TxApplyProcessing: tw.Tx.UnsafeMeta,
	}

	bz, _ := txResultMeta.MarshalBinary()

	return bz
}

func (tw *IndexerTransactionExtractor) GetTransaction() *Transaction {
	return &Transaction{
		Hash:             tw.GetTransactionHash(),
		Status:           tw.GetStatus(),
		Ledger:           tw.GetLedgerSequence(),
		ApplicationOrder: tw.GetApplicationOrder(),
		EnvelopeXdr:      tw.GetEnvelopeXdr(),   // xdr.TransactionEnvelope
		ResultXdr:        tw.GetResultXdr(),     // xdr.TransactionResultPair
		ResultMetaXdr:    tw.GetResultMetaXdr(), //xdr.TransactionResultMeta
		SourceAddress:    tw.Tx.Envelope.SourceAccount().ToAccountId().Address(),
		TransactionTime:  tw.Time,
	}
}
