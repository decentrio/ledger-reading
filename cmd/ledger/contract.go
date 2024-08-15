package main

import (
	"github.com/decentrio/ledger-reading/manager"
	"github.com/spf13/cobra"
)

// NewContractRead
func NewContractRead() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "contract [contract_id]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			contractId := args[0]
			m := manager.DefaultNewManager()
			m.Reader.ContractReading(contractId)
			return nil
		},
	}

	return cmd
}
