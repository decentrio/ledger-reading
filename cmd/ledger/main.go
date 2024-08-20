package main

import (
	"os"
	"path/filepath"

	"github.com/decentrio/ledger-reading/lib/cli"
)

func main() {
	rootCmd.AddCommand(NewRunNodeCmd())
	rootCmd.AddCommand(NewContractRead())
	rootCmd.AddCommand(NewContractTxsRead())

	cmd := cli.PrepareBaseCmd(rootCmd, "CMT", os.ExpandEnv(filepath.Join("$HOME", DefaultCometDir)))
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
