package main

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewCloneCommand() *cobra.Command {
	cloneCommand := &cobra.Command{
		Use:   "clone SRC_NETWORK DEST_NETWORK",
		Short: "Copy routes between networks",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if args[0] == "" {
				return errors.New("source can not be empty")
			}
			if args[1] == "" {
				return errors.New("destination can not be empty")
			}

			fmt.Println("Source network:", args[0])
			fmt.Println("Destination network:", args[1])

			return nil
		},
	}

	// cloneCommand.Flags().String("from", "", "Network ID to copy from (required)")
	// cloneCommand.Flags().String("to", "", "Network ID to copy to (required)")

	// cloneCommand.MarkFlagRequired("from")
	// cloneCommand.MarkFlagRequired("to")
	// cloneCommand.MarkFlagsRequiredTogether("from", "to")

	return cloneCommand
}
