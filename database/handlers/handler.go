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
