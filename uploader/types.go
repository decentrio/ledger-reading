package uploader

import (
	"github.com/decentrio/ledger-reading/importer"
	"github.com/stellar/go/xdr"
)

const (
	Buy  = "buy"
	Sell = "sell"
)

const (
	UsdcTokenName           = "USDC-GA5ZSEJYB37JRC5AVCIA5MOP4RHTM335X2KGX3IHOJAPP5RE34K4KZVN"
	MultiHopSorobanContract = "CCLZRD4E72T7JCZCN3P7KNPYNXFYKQCL64ECLX7WP5GNVYPYJGU2IO2G"
	BaseCurrencyKey         = 1
	TargetCurrencyKey       = 2
)

type ContractData struct {
	ContractId string
	AccountId  string
	Key        []xdr.ScVal
	Value      []xdr.ScVal
}

type WasmContractEvent struct {
	Id         string
	ContractId string
	TxHash     string
	EventBody  []xdr.ContractEventBody
}

type InvokeTransaction struct {
	Hash         string `json:"hash,omitempty"`
	ContractId   string `json:"contract_id,omitempty"`
	FunctionType string `json:"func_type,omitempty"`
	FunctionName string `json:"function_name,omitempty"`
	Args         []byte `json:"args,omitempty"`
}

type Transaction struct {
	Hash             string `json:"hash,omitempty"`
	Status           string `json:"status,omitempty"`
	Ledger           uint32 `json:"ledger,omitempty"`
	ApplicationOrder uint32 `json:"application_order,omitempty"`
	EnvelopeXdr      []byte `json:"envelope_xdr,omitempty"`
	ResultXdr        []byte `json:"result_xdr,omitempty"`
	ResultMetaXdr    []byte `json:"result_meta_xdr,omitempty"`
	SourceAddress    string `json:"source_address,omitempty"`
	TransactionTime  uint64 `json:"transaction_time,omitempty"`
}

type LegacyTradeInformation struct {
	OfferTokenSorobanContract string
	OfferAmount               uint64
	AskedTokenSorobanContract string
	ReturnAmount              uint64
	Price                     float64
}

type TradeInformation struct {
	TradeType    string
	Ticker       importer.Ticker
	BaseVolume   uint64
	TargetVolume uint64
	Price        float64
}

type ProvideLiquidityInformation struct {
	Sender              string          `json:"sender,omitempty"`
	Ticker              importer.Ticker `json:"ticker,omitempty"`
	BaseVolume          uint64          `json:"base_volume,omitempty"`
	TargetVolume        uint64          `json:"target_volume,omitempty"`
	ShareLiquidity      uint64          `json:"share_liquidity,omitempty"`
	BaseLiquidity       uint64          `json:"base_liquidity,omitempty"`
	TargetLiquidity     uint64          `json:"target_liquidity,omitempty"`
	TotalLiquidityInUsd uint64          `json:"total_liquidity_in_usd,omitempty"`
}

type Ticker struct {
	TickerId       string `json:"ticker_id,omitempty"`        // PHO_USDC
	BaseCurrency   string `json:"base_currency,omitempty"`    // PHO
	TargetCurrency string `json:"target_currency,omitempty"`  // USDC
	PoolId         string `json:"pool_id,omitempty"`          // "CAZ6W4WHVGQBGURYTUOLCUOOHW6VQGAAPSPCD72VEDZMBBPY7H43AYEC"
	LastPrice      string `json:"last_price,omitempty"`       // Last price trade
	BaseVolume     uint64 `json:"base_volume,omitempty"`      // base currency trade volume (24h)
	TargetVolume   uint64 `json:"target_volume,omitempty"`    // target currency trade volume (24h)
	LiquidityInUsd uint64 `json:"liquidity_in_usd,omitempty"` // liquidity in usd
}
