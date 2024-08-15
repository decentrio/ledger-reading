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

type ContractsData struct {
	Id            string `json:"id,omitempty"`
	ContractId    string `json:"contract_id,omitempty"`
	AccountId     string `json:"account_id,omitempty"`
	TxHash        string `json:"tx_hash,omitempty"`
	Ledger        uint32 `json:"ledger,omitempty"`
	EntryType     string `json:"entry_type,omitempty"`
	KeyXdr        []byte `json:"key_xdr,omitempty"`
	ValueXdr      []byte `json:"value_xdr,omitempty"`
	Durability    int32  `json:"durability,omitempty"`
	IsNewest      bool   `json:"is_newest,omitempty"`
	UpdatedLedger uint32 `json:"updated_ledger,omitempty"` // previous updated ledger (TODO: we should correct the name here)
}
