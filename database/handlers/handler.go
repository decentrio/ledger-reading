package handlers

import (
	"github.com/decentrio/ledger-reading/database/models"
)

func (h *DBHandler) TransactionsAtLedgerSeq(ledger int32) (data []*models.Transaction, err error) {
	err = h.db.Table("transactions").Where("seq = ?", ledger).Find(&data).Error
	if err != nil {
		return []*models.Transaction{}, err
	}

	return data, nil
}

func (h *DBHandler) TransactionByHash(hash string) (data *models.Transaction, err error) {
	err = h.db.Table("transactions").Where("hash = ?", hash).First(&data).Error
	if err != nil {
		return &models.Transaction{}, err
	}

	return data, nil
}

func (h *DBHandler) ContractData(contractId string) (data []*models.ContractsData, err error) {
	err = h.db.Table("contracts_data").
	Where("contract_id = ?", contractId).
	Find(&data).Error
	if err != nil {
		return []*models.ContractsData{}, err
	}

	return data, nil
}