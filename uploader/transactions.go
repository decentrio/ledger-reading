package uploader

import (
	"fmt"

	"github.com/decentrio/converter/converter"
	"github.com/stellar/go/ingest"
	"github.com/stellar/go/xdr"
)

const (
	SUCCESS = "success"
	FAILED  = "failed"
)

type PhoenixTransactionExtractor struct {
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

type PhoenixActionTx struct {
	Type ActionType
	Tx   *PhoenixTransactionExtractor
}

func NewPhoenixTransactionExtractor(tx ingest.LedgerTransaction, ledgerSeq uint32, processedUnixTime uint64) *PhoenixActionTx {
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

	txExt := &PhoenixTransactionExtractor{
		LedgerSequence: ledgerSeq,
		Tx:             tx,
		Ops:            ops,
		Time:           processedUnixTime,
	}

	actionType, err := txExt.GetPhoenixActionType()
	if err != nil {
		return nil
	}

	return &PhoenixActionTx{
		Type: actionType,
		Tx:   txExt,
	}
}

func (tx *PhoenixTransactionExtractor) GetPhoenixActionType() (ActionType, error) {
	invoke, isIvkFunc := tx.IsInvokeHostFunctionTx()
	if !isIvkFunc {
		return Unkown, fmt.Errorf("not a invoke host function")
	}

	var argsXdr xdr.InvokeContractArgs
	argsXdr.UnmarshalBinary(invoke.Args)
	method := string(argsXdr.FunctionName)

	if method == "swap" {
		fmt.Printf("============\n\nswap\n\n============")
		return Swap, nil
	}

	if method == "provide_liquidity" {
		fmt.Printf("============\n\nprovide_liquidity\n\n============")
		return ProvideLiquidity, nil
	}

	return Unkown, fmt.Errorf("unkown method")
}

func (tx *PhoenixTransactionExtractor) IsInvokeHostFunctionTx() (InvokeTransaction, bool) {
	var invokeFuncTx InvokeTransaction
	var isInvokeFuncTx bool

	ops := tx.Tx.Envelope.Operations()
	for _, op := range ops {
		if op.Body.Type == xdr.OperationTypeInvokeHostFunction {
			ihfOp := op.Body.MustInvokeHostFunctionOp()
			switch ihfOp.HostFunction.Type {
			case xdr.HostFunctionTypeHostFunctionTypeInvokeContract:

				ic := ihfOp.HostFunction.MustInvokeContract()
				ca, err := converter.ConvertScAddress(ic.ContractAddress)
				if err != nil {
					continue
				}

				fn := string(ic.FunctionName)

				args, err := ic.MarshalBinary()
				if err != nil {
					continue
				}

				invokeFuncTx.Hash = tx.Tx.Result.TransactionHash.HexString()
				invokeFuncTx.ContractId = *ca.ContractId
				invokeFuncTx.FunctionType = "invoke_host_function"
				invokeFuncTx.FunctionName = fn
				invokeFuncTx.Args = args

				isInvokeFuncTx = true

				break
			case xdr.HostFunctionTypeHostFunctionTypeCreateContract:
				// we do not care about this type
				continue

			case xdr.HostFunctionTypeHostFunctionTypeUploadContractWasm:
				// we do not care about this type
				continue
			}

		}
	}

	return invokeFuncTx, isInvokeFuncTx
}

func (tw *PhoenixTransactionExtractor) GetTransactionHash() string {
	return tw.Tx.Result.TransactionHash.HexString()
}

func (tw *PhoenixTransactionExtractor) GetStatus() string {
	if tw.Tx.Result.Successful() {
		return SUCCESS
	}

	return FAILED
}

func (tw *PhoenixTransactionExtractor) GetLedgerSequence() uint32 {
	return tw.LedgerSequence
}

func (tw *PhoenixTransactionExtractor) GetApplicationOrder() uint32 {
	return tw.Tx.Index
}

func (tw *PhoenixTransactionExtractor) GetEnvelopeXdr() []byte {
	bz, _ := tw.Tx.Envelope.MarshalBinary()
	return bz
}

func (tw *PhoenixTransactionExtractor) GetResultXdr() []byte {
	bz, _ := tw.Tx.Result.MarshalBinary()
	return bz
}

func (tw *PhoenixTransactionExtractor) GetResultMetaXdr() []byte {
	txResultMeta := xdr.TransactionResultMeta{
		Result:            tw.Tx.Result,
		FeeProcessing:     tw.Tx.FeeChanges,
		TxApplyProcessing: tw.Tx.UnsafeMeta,
	}

	bz, _ := txResultMeta.MarshalBinary()

	return bz
}

func (tw *PhoenixTransactionExtractor) GetTransaction() *Transaction {
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
