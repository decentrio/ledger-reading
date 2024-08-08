package models

type Tokens struct {
	Symbol          string  `json:"symbol,omitempty"`
	TokenName       string  `json:"token_name,omitempty"`
	SorobanContract string  `json:"soroban_contract,omitempty"`
	Decimal         int     `json:"decimal,omitempty"`
	PriceInUsd      float64 `json:"price_in_usd,omitempty"`
}

type Tickers struct {
	TickerId        string  `json:"ticker_id,omitempty"`        // PHO_USDC
	BaseCurrency    string  `json:"base_currency,omitempty"`    // PHO
	TargetCurrency  string  `json:"target_currency,omitempty"`  // USDC
	PoolId          string  `json:"pool_id,omitempty"`          // "CAZ6W4WHVGQBGURYTUOLCUOOHW6VQGAAPSPCD72VEDZMBBPY7H43AYEC"
	LastPrice       float64 `json:"last_price,omitempty"`       // Last price trade
	BaseVolume      uint64  `json:"base_volume,omitempty"`      // base currency trade volume (24h)
	TargetVolume    uint64  `json:"target_volume,omitempty"`    // target currency trade volume (24h)
	High            float64 `json:"high,omitempty"`             // highest price in 24h
	Low             float64 `json:"low,omitempty"`              // lowest price in 24h
	ShareLiquidity  uint64  `json:"share_liquidity,omitempty"`  // current share liquidity
	BaseLiquidity   uint64  `json:"base_liquidity,omitempty"`   // current base liquidity
	TargetLiquidity uint64  `json:"target_liquidity,omitempty"` // current target liquidity
	LiquidityInUsd  uint64  `json:"liquidity_in_usd,omitempty"` // liquidity in usd
}

type HistoricalTrades struct {
	TradeId        uint32  `json:"trade_id,omitempty"`        // A unique ID associated with the trade for the currency pair transaction
	Price          float64 `json:"price,omitempty"`           // Transaction price of base asset in target currency
	TickerId       string  `json:"ticker_id,omitempty"`       // Ticker ID
	BaseVolume     uint64  `json:"base_volume,omitempty"`     // volume trade of base currency (float)
	TargetVolume   uint64  `json:"target_volume,omitempty"`   // volume trade of target currency (float)
	TradeTimestamp uint64  `json:"trade_timestamp,omitempty"` // time stamp of trade
	TradeType      string  `json:"trade_type,omitempty"`      // buy/sell
	TxHash         string  `json:"tx_hash,omitempty"`
	Maker          string  `json:"maker,omitempty"`
}

type Activities struct {
	Address        string
	ActionType     string
	BaseCurrency   string
	BaseVolume     uint64
	TargetCurrency string
	TargetVolume   uint64
	Timestamp      uint64
}
