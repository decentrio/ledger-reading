package uploader

import (
	"github.com/decentrio/ledger-reading/importer"
)

type TransactionExtractor interface {
	IsInvokeHostFunctionTx() (InvokeTransaction, error)
	GetContractEvents() ([]WasmContractEvent, error)
	GetContractDataEntry() []ContractData
	GetHistoricalLegacyTradeInformation() (LegacyTradeInformation, error)
	GetHistoricalTradeInformation() (TradeInformation, error)
	GetTradeTicker(tickerList map[string]importer.Ticker) (Ticker, error)
}
