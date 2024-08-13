package models

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
