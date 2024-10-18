package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewWhatchCommand() *cobra.Command {
	whatchCommand := &cobra.Command{
		Use:   "whatch [OPTIONS] NETWORK [NETWORKS...]",
		Short: "Run updates periodically",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("whatch: not implemented")
		},
	}

	whatchCommand.Flags().String("interval", "8", "Time interval between update in hours (required)")
	whatchCommand.MarkFlagRequired("interval")

	return whatchCommand
}
