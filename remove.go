package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewRemoveCommand() *cobra.Command {
	removeCommand := &cobra.Command{
		Use:   "rm [OPTIONS] NETWORK DOMAIN [DOMAIN...]",
		Short: "Remove domains. If --all is specified, remove all domains",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("remove: not implemented")
		},
	}

	removeCommand.Flags().StringP("group", "g", "default", "Group from which to remove domain")
	removeCommand.Flags().Bool("all", false, "Remove all domains")

	return removeCommand
}
