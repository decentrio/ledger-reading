package uploader

import (
	"fmt"

	"github.com/stellar/go/ingest"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
	"golang.org/x/exp/maps"
)

const (
	// Implemented
	EventTypeTransfer = "transfer"
	EventTypeMint     = "mint"
	EventTypeClawback = "clawback"
	EventTypeBurn     = "burn"
	// TODO: Not implemented
	EventTypeIncrAllow
	EventTypeDecrAllow
	EventTypeSetAuthorized
	EventTypeSetAdmin
)

var (
	STELLAR_ASSET_CONTRACT_TOPICS = map[xdr.ScSymbol]string{
		xdr.ScSymbol("transfer"): EventTypeTransfer,
		xdr.ScSymbol("mint"):     EventTypeMint,
		xdr.ScSymbol("clawback"): EventTypeClawback,
		xdr.ScSymbol("burn"):     EventTypeBurn,
	}
)

func (tx *IndexerTransactionExtractor) GetContractEvents() ([]WasmContractEvent, error) {
	wasmContractevents := make(map[string]WasmContractEvent)
	for _, op := range tx.Ops {
		var order = uint32(1)
		if op.OperationType() == xdr.OperationTypeInvokeHostFunction {
			diagnosticEvents, innerErr := tx.Tx.GetDiagnosticEvents()
			if innerErr != nil {
				return nil, innerErr
			}
			evts := filterEvents(diagnosticEvents)

			for _, evt := range evts {
				isAssetEvent := isStellarAssetContractEvent(evt)
				if !isAssetEvent {
					wasmEvent, err := getWasmContractEvents(tx.Tx, evt, op.ID(), &order)
					if err != nil {
						continue
					}

					contractEvent, found := wasmContractevents[wasmEvent.ContractId]
					if found {
						contractEvent.EventBody = append(contractEvent.EventBody, wasmEvent.EventBody...)
						wasmContractevents[wasmEvent.ContractId] = contractEvent
					} else {
						wasmContractevents[wasmEvent.ContractId] = wasmEvent
					}
				}
			}
		}
	}

	contractEvents := maps.Values(wasmContractevents)
	return contractEvents, nil
}

func getWasmContractEvents(tx ingest.LedgerTransaction, event xdr.ContractEvent, id int64, order *uint32) (WasmContractEvent, error) {
	contractId, err := strkey.Encode(strkey.VersionByteContract, event.ContractId[:])
	if err != nil {
		return WasmContractEvent{}, err
	}

	evt := WasmContractEvent{
		Id:         fmt.Sprintf("%019d-%010d", id, *order), // ID should be combine from operation ID and event index
		ContractId: contractId,
		TxHash:     tx.Result.TransactionHash.HexString(),
		EventBody:  []xdr.ContractEventBody{event.Body},
	}
	*order++

	return evt, nil
}

func isStellarAssetContractEvent(event xdr.ContractEvent) bool {
	if event.Type != xdr.ContractEventTypeContract || event.ContractId == nil || event.Body.V != 0 {
		return false
	}

	topics := event.Body.V0.Topics

	// No relevant SAC events have <= 2 topics
	if len(topics) <= 2 {
		return false
	}

	fn, ok := topics[0].GetSym()
	if !ok {
		return false
	}

	if _, found := STELLAR_ASSET_CONTRACT_TOPICS[fn]; !found {
		return false
	}

	return true
}

func filterEvents(diagnosticEvents []xdr.DiagnosticEvent) []xdr.ContractEvent {
	var filtered []xdr.ContractEvent
	for _, diagnosticEvent := range diagnosticEvents {
		if !diagnosticEvent.InSuccessfulContractCall || diagnosticEvent.Event.Type != xdr.ContractEventTypeContract {
			continue
		}
		filtered = append(filtered, diagnosticEvent.Event)
	}
	return filtered
}
