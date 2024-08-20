package reader

import "github.com/stellar/go/xdr"

type InvokeTransaction struct {
	Hash         string `json:"hash,omitempty"`
	ContractId   string `json:"contract_id,omitempty"`
	FunctionType string `json:"func_type,omitempty"`
	FunctionName string `json:"function_name,omitempty"`
	Args         []byte `json:"args,omitempty"`
}

type ContractData struct {
	ContractId string
	AccountId  string
	Key        []xdr.ScVal
	Value      []xdr.ScVal
}