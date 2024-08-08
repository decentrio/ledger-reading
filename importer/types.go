package importer

const (
	TokenListUrl  = "https://raw.githubusercontent.com/decentrio/token-list/main/token_lists.json"
	TickerListUrl = "https://raw.githubusercontent.com/decentrio/token-list/main/ticker_lists.json"
)

func GetTokenPair(base string, target string) TokenPair {
	return TokenPair{
		BaseCurrency:   base,
		TargetCurrency: target,
	}
}

type TokenPair struct {
	BaseCurrency   string `json:"base_currency,omitempty"`   // token_a
	TargetCurrency string `json:"target_currency,omitempty"` // token_b
}

type Ticker struct {
	TickerId       string `json:"ticker_id,omitempty"`
	BaseCurrency   string `json:"base_currency,omitempty"`   // token_a
	TargetCurrency string `json:"target_currency,omitempty"` // token_b
	PoolContract   string `json:"pool_contract,omitempty"`
}

type Token struct {
	Symbol          string `json:"symbol,omitempty"`
	Token           string `json:"token,omitempty"`
	SorobanContract string `json:"soroban_contract,omitempty"`
	Decimals        uint32 `json:"decimals,omitempty"`
}
