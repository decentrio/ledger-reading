package uploader

type TransactionExtractor interface {
	IsInvokeHostFunctionTx() (InvokeTransaction, error)
	GetContractEvents() ([]WasmContractEvent, error)
	GetContractDataEntry() []ContractData
	GetHistoricalLegacyTradeInformation() (LegacyTradeInformation, error)
	GetHistoricalTradeInformation() (TradeInformation, error)
	GetTradeTicker(tickerList map[string]ITicker) (Ticker, error)
}
