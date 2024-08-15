package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/decentrio/ledger-reading/manager"
	"github.com/decentrio/ledger-reading/lib/cli"
	"github.com/spf13/cobra"
)

var (
	DefaultCometDir = ".ledger"
)

var rootCmd = &cobra.Command{
	Use: "ledger",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

// NewRunNodeCmd returns the command that allows the CLI to start a node.
// It can be used with a custom PrivValidator and in-process ABCI application.
func NewRunNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"node", "run"},
		RunE: func(cmd *cobra.Command, args []string) error {
			from, to := ParseConfig(cmd)

			m := manager.NewManagerFromTo(from, to)

			if err := m.Start(); err != nil {
				return fmt.Errorf("failed to start node: %w", err)
			}

			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)

			go func() {
				for range c {
					if m.IsRunning() {
						if err := m.Stop(); err != nil {
							fmt.Printf(err.Error())
						}
					}
					os.Exit(0)
				}
			}()

			// Run forever.
			select {}
		},
	}

	return cmd
}

func ParseConfig(cmd *cobra.Command) (int32, int32) {
	fromLedger, err := cmd.Flags().GetInt32(cli.FromLedger)
	if err != nil {
		fromLedger = 53012912
	}

	toLedger, err := cmd.Flags().GetInt32(cli.ToLedger)
	if err != nil {
		toLedger = fromLedger +1
	}
	
	return fromLedger, toLedger
}