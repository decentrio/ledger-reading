package reader

import (
	"fmt"
	"time"

	"github.com/decentrio/ledger-reading/database/models"
	"github.com/decentrio/xdr-converter/converter"
	"github.com/stellar/go/ingest"
	"github.com/stellar/go/xdr"
	"github.com/stellar/go/strkey"
)

func (r *Reader) LedgerProcessing() {
	fmt.Println("from ledger:", r.FromSeq)
	fmt.Println("to ledger:", r.ToSeq)
	for i := r.FromSeq; i <= r.ToSeq; i++ {
		var contractTxs []ingest.LedgerTransaction
		txs, err := r.db.TransactionsAtLedgerSeq(i)
		if err != nil {
			panic(err)
		}
		fmt.Println("ledger:", i)

		fmt.Println("count txs:", len(txs))
		for _, tx := range txs {
			orgTx := convertToOriginTx(tx)
			ivkFunc, isIvk := IsInvokeHostFunctionTx(&orgTx)

			if isIvk {

				contractTxs = append(contractTxs, orgTx)
				var argsXdr xdr.InvokeContractArgs
				argsXdr.UnmarshalBinary(ivkFunc.Args)
				contractId := ""
				if argsXdr.ContractAddress.ContractId != nil {
					contractId, err = strkey.Encode(strkey.VersionByteContract, argsXdr.ContractAddress.ContractId[:])
					if err != nil {
						continue
					}
				}
				method := string(argsXdr.FunctionName)
				fmt.Println(method)
				if method == "swap" {
					fmt.Println(contractId)
				}
			}

		}
		fmt.Println("count contract txs:",len(contractTxs))
		time.Sleep(time.Millisecond * 500)
	}
}

func convertToOriginTx(tx *models.Transaction) ingest.LedgerTransaction {
	var envelop xdr.TransactionEnvelope
	err := envelop.UnmarshalBinary(tx.EnvelopeXdr)
	if err != nil {
		panic(err)
	}
	var resultMeta xdr.TransactionResultMeta
	err = resultMeta.UnmarshalBinary(tx.ResultMetaXdr)
	if err != nil {
		panic(err)
	}

	return ingest.LedgerTransaction{
		Index:         tx.ApplicationOrder,
		Envelope:      envelop,
		Result:        resultMeta.Result,
		FeeChanges:    resultMeta.FeeProcessing,
		UnsafeMeta:    resultMeta.TxApplyProcessing,
		LedgerVersion: 1,
	}
}


func IsInvokeHostFunctionTx(tx *ingest.LedgerTransaction) (InvokeTransaction, bool) {
	var invokeFuncTx InvokeTransaction
	var isInvokeFuncTx bool

	ops := tx.Envelope.Operations()
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

				invokeFuncTx.Hash = tx.Result.TransactionHash.HexString()
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