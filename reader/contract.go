package reader

import (
	"fmt"
	"golang.org/x/exp/maps"

	"github.com/stellar/go/ingest"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
)

func (r *Reader) ContractReading(contractId string) {
	entries, err := r.db.ContractData(contractId)

	if err != nil {
		panic(err)
	}
	fmt.Println(len(entries))
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

func (r *Reader) ContractReadingTxs(contractId string) {
	entries, err := r.db.ContractData(contractId)

	if err != nil {
		panic(err)
	}
	for _, entry := range entries {
		tx, err := r.db.TransactionByHash(entry.TxHash)
		if err != nil {
			panic(err)
		}
		orgTx := ConvertToOriginTx(tx)
		ivkFuncs, isIvk := IsInvokeHostFunctionTx(&orgTx)
		fmt.Println(len(ivkFuncs))
		if isIvk {
			for _, ivk := range ivkFuncs {
				var argsXdr xdr.InvokeContractArgs
				argsXdr.UnmarshalBinary(ivk.Args)
				// contractId := ""
				// if argsXdr.ContractAddress.ContractId != nil {
				// 	contractId, err = strkey.Encode(strkey.VersionByteContract, argsXdr.ContractAddress.ContractId[:])
				// 	if err != nil {
				// 		continue
				// 	}
				// }
				method := string(argsXdr.FunctionName)
				fmt.Println(method)
			}

			contractData := GetContractDataEntry(orgTx)

			for _, entry := range contractData {
				fmt.Println("==============================")
				fmt.Println(entry.ContractId)
				fmt.Println("/////==============================/////")
				fmt.Println(entry.Key)
				fmt.Println("==============================")

			}
		}
	}
}

func GetContractDataEntry(tx ingest.LedgerTransaction) []ContractData {
	v3 := tx.UnsafeMeta.V3
	if v3 == nil {
		return nil
	}

	contract := make(map[string]ContractData)
	for _, op := range v3.Operations {
		for _, change := range op.Changes {
			entry, entryType, found := ContractDataEntry(change)
			// continue with "state" because we don't want to store this entry
			if entryType == "state" {
				continue
			}
			if found {
				var contractId string
				var err error
				if entry.Contract.ContractId != nil {
					contractId, err = strkey.Encode(strkey.VersionByteContract, entry.Contract.ContractId[:])
					if err != nil {
						continue
					}
				}

				contractData, found := contract[contractId]
				if found {
					contractData.Key = append(contractData.Key, entry.Key)
					contractData.Value = append(contractData.Value, entry.Val)
					contract[contractId] = contractData
				} else {
					var accountId string
					if entry.Contract.AccountId != nil {
						accountId, err = entry.Contract.AccountId.GetAddress()
						if err != nil {
							continue
						}
					}

					cd := ContractData{
						ContractId: contractId,
						AccountId:  accountId,
						Key:        []xdr.ScVal{entry.Key},
						Value:      []xdr.ScVal{entry.Val},
					}

					contract[contractId] = cd
				}
			}
		}
	}

	contractData := maps.Values(contract)
	return contractData
}

func ReadMapValue(xdrMap *xdr.ScMap, key xdr.ScVal) xdr.ScVal {
	for _, entry := range *xdrMap {
		if entry.Key.Equals(key) {
			return entry.Val
		}
	}
	return xdr.ScVal{}
}

func ContractDataEntry(c xdr.LedgerEntryChange) (xdr.ContractDataEntry, string, bool) {
	var result xdr.ContractDataEntry

	switch c.Type {
	case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
		created := *c.Created
		if created.Data.ContractData != nil {
			result = *created.Data.ContractData
			return result, "created", true
		}
	case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
		updated := *c.Updated
		if updated.Data.ContractData != nil {
			result = *updated.Data.ContractData
			return result, "updated", true
		}
	case xdr.LedgerEntryChangeTypeLedgerEntryRemoved:
		ledgerKey := c.Removed
		if ledgerKey.ContractData != nil {
			result.Contract = ledgerKey.ContractData.Contract
			result.Key = ledgerKey.ContractData.Key
			result.Durability = ledgerKey.ContractData.Durability
			return result, "removed", true
		}
	case xdr.LedgerEntryChangeTypeLedgerEntryState:
		state := *c.State
		if state.Data.ContractData != nil {
			result = *state.Data.ContractData
			return result, "state", true
		}

	}
	return result, "", false
}
