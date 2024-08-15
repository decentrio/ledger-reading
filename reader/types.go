package reader

type InvokeTransaction struct {
	Hash         string `json:"hash,omitempty"`
	ContractId   string `json:"contract_id,omitempty"`
	FunctionType string `json:"func_type,omitempty"`
	FunctionName string `json:"function_name,omitempty"`
	Args         []byte `json:"args,omitempty"`
}
