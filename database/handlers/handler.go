package handlers

import (
	"github.com/decentrio/ledger-reading/database/models"
)

func (h *DBHandler) TransactionsAtLedgerSeq(ledger int32) (data []*models.Transaction, err error) {
	err = h.db.Table("transactions").Where("ledger = ?", ledger).Find(&data).Error
	if err != nil {
		return []*models.Transaction{}, err
	}

	return data, nil
}

func (h *DBHandler) ContractData(contractId string) (data []*models.ContractsData, err error) {
	err = h.db.Table("contracts_data").
	Where("contract_id = ?", contractId).
	Where("is_newest = ?", true).
	Find(&data).Error
	if err != nil {
		return []*models.ContractsData{}, err
	}

	return data, nil
}