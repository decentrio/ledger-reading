package uploader

import (
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
	"golang.org/x/exp/maps"
)

const (
	ShareTokenKey  = uint32(0)
	BaseTokenKey   = uint32(1)
	TargetTokenKey = uint32(2)
)

func (tw *PhoenixTransactionExtractor) GetPoolLiquidity(data ContractData) (uint64, uint64, uint64, error) {
	var shareLiquidity uint64
	var baseLiquidity uint64
	var targetLiquidity uint64

	for index, keyXdr := range data.Key {
		key := uint32(keyXdr.MustU32())
		val := data.Value[index].MustI128()

		switch key {
		case ShareTokenKey:
			shareLiquidity = uint64(val.Lo)
		case BaseTokenKey:
			baseLiquidity = uint64(val.Lo)
		case TargetTokenKey:
			targetLiquidity = uint64(val.Lo)
		default:
		}
	}
	return shareLiquidity, baseLiquidity, targetLiquidity, nil
}

func (tw *PhoenixTransactionExtractor) GetContractDataEntry() []ContractData {
	v3 := tw.Tx.UnsafeMeta.V3
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
