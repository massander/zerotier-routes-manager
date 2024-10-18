package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	ztroutesCommand := &cobra.Command{
		Use:     "ztroutes [OPTIONS] COMMAND [ARG...]",
		Short:   "ztroutes is a CLI tool for configuring ZeroTier Managed Routes",
		Version: "v0.0.2",
	}

	ztroutesCommand.PersistentFlags().StringP("config", "c", ZTROUTES_CONFIG, "The location of your routes configuration files ")
	ztroutesCommand.PersistentFlags().Bool("debug", false, "If provided command will not push changes to ZeroTier Network Controller and change only local config and print result")
	ztroutesCommand.PersistentFlags().StringP("token", "t", "", "ZeroTier Auth token from contoller (required)")

	// ztroutesCommand.MarkPersistentFlagRequired("token")

	ztroutesCommand.AddCommand(
		NewLookupCommand(),
		NewAddCommand(),
		NewRemoveCommand(),
		NewCloneCommand(),
	)

	if err := ztroutesCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
