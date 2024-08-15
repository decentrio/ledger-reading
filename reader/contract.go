package reader

import (
	"fmt"

	"github.com/stellar/go/xdr"
)

func (r *Reader) ContractReading(contractId string) {
	entries, err := r.db.ContractData(contractId)

	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		var keySc xdr.ScVal
		var valSc xdr.ScVal

		keySc.UnmarshalBinary(entry.KeyXdr)
		valSc.UnmarshalBinary(entry.ValueXdr)
		if keySc.Address != nil {
			fmt.Println(keySc.Address.AccountId.Address())
			fmt.Println(valSc.Type)
			fmt.Println(valSc.MustMap())
		}
	}
}

