package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewAddCommand() *cobra.Command {
	addCommand := &cobra.Command{
		Use:   "add [OPTIONS] NETWORK DOMAIN [DOMAIN...] ",
		Short: "Add domain or IP adresses to routes configuration",
		Args:  cobra.RangeArgs(1, 5),
		RunE: func(cmd *cobra.Command, args []string) error {
			network := args[0]
			fmt.Println("Network:", network)
			fmt.Println("Values:", args[1:])
			return nil
		},
	}

	addCommand.Flags().StringP("group", "g", "default", "Group to which to add domain")

	return addCommand
}
