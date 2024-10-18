package main

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func NewLookupCommand() *cobra.Command {
	lookupCommand := &cobra.Command{
		Use:   "lookup NETWORK",
		Short: "Look up IP addresses and update ZeroTier routes",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()

			configDir, err := flags.GetString("config")
			if err != nil {
				return fmt.Errorf("flags.GetConfig: %w", err)
			}

			// if path == "" {
			// 	return errors.New("config path can not be empty")
			// }

			debug, err := flags.GetBool("debug")
			if err != nil {
				return err
			}

			token, err := flags.GetString("token")
			// Return error only if not in debug mode
			if !debug {
				if err != nil {
					return err
				}

				if token == "" {
					return errors.New("flag not provided: token")
				}
			}

			network := args[0]
			if network == "" {
				return errors.New("argument can not be empty: NETWORK")
			}

			path := fmt.Sprintf("%s/%s.routes.json", configDir, network)

			config, err := loadConfig(path)
			cobra.CheckErr(err)

			err = lookupIPs(&config)
			cobra.CheckErr(err)

			err = saveConfig(config, path)
			cobra.CheckErr(err)

			if !debug {
				err := syncRoutes(config, network, token)
				cobra.CheckErr(err)
			} else {
				configBytes, err := json.MarshalIndent(config, "", "  ")
				cobra.CheckErr(err)

				fmt.Println(string(configBytes))
			}

			return nil
		},
	}

	return lookupCommand
}

func lookupIPs(config *Config) error {
	const operation = "lookupIPs"
	for _, group := range *config {
		if !group.IsEnabled {
			continue
		}

		for i, route := range group.Routes {
			ipsv4 := make([]string, 0)

			ips, err := net.LookupIP(route.Domain)
			if err != nil {
				return fmt.Errorf("%s: %w", operation, err)
			}

			for _, ip := range ips {
				// Skip IPv6
				if ip.DefaultMask() == nil {
					fmt.Println(ip, "is IPv6")
					continue
				}

				ipsv4 = append(ipsv4, ip.String())
			}

			group.Routes[i].IPs = ipsv4
		}
	}

	return nil
}
