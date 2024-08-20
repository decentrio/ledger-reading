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
			valMap, ok := valSc.GetMap()
			if !ok {
				fmt.Println("This is not a map")
			} else {
				var stakeKeySym xdr.ScSymbol = "total_stake"
				stakeKey, _ := xdr.NewScVal(xdr.ScValTypeScvSymbol, stakeKeySym)
				stakeValue := ReadMapValue(valMap, stakeKey)

				stakeI128, ok := stakeValue.GetI128()
				if ok {
					fmt.Println("total stake value hi:", stakeI128.Hi)
					fmt.Println("total stake value lo:", stakeI128.Lo)
				}
			}
		}
	}
}

func ReadMapValue(xdrMap *xdr.ScMap, key xdr.ScVal) xdr.ScVal {
	for _, entry := range *xdrMap {
		if entry.Key.Equals(key) {
			return entry.Val
		}
	}
	return xdr.ScVal{}
}
