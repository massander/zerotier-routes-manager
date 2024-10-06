package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/hjson/hjson-go/v4"
	"github.com/spf13/cobra"
)

type App struct {
	Name string `json:"name"`
	// List of domains to route.
	Domains []string `json:"domains"`
	IPs     []string `json:"ips,omitempty"`
	// Zero Tier routes config.
	Routes []ZTRoute `json:"routes"`
}

type ZTNetwork struct {
	// Device (linux server) through which to route traffic.
	NetworkID string `json:"networkId"`
	ExitNode  string `json:"exitNode"`
	LAN       string `json:"lan"`
	Apps      []App  `json:"apps"`
}

type ZTRoute struct {
	Target string `json:"target"`
	Via    string `json:"via,omitempty"`
}

type ZTNetworkConfig struct {
	Routes []ZTRoute `json:"routes"`
}

type ZTConfig struct {
	Config ZTNetworkConfig `json:"config"`
}

const ZT_URL = "https://api.zerotier.com/api/v1"

func main() {
	var configPath string
	var debug bool

	var rootCommand = &cobra.Command{
		Use:     "zt-routes",
		Short:   "zt-routes is a CLI tool for managing ZeroTier Managed Routes ",
		Version: "v0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {

			ztNet, err := loadConfig(configPath)
			if err != nil {
				return err
			}

			if ztNet.NetworkID == "" {
				return fmt.Errorf("Network ID is required")
			}

			if ztNet.ExitNode == "" {
				return fmt.Errorf("Exit Node address is required")
			}

			if ztNet.LAN == "" {
				return fmt.Errorf("LAN is required")
			}

			for i, app := range ztNet.Apps {
				ztRoutes := make([]ZTRoute, 0)
				for _, domain := range app.Domains {

					ips, err := net.LookupIP(domain)
					if err != nil {
						return fmt.Errorf("lookupIPs: %w", err)
					}

					for _, ip := range ips {
						if ip.DefaultMask() == nil {
							fmt.Println(ip, "is IPv6. Skipping")
							continue
						}
						ztRoutes = append(ztRoutes, ZTRoute{
							Target: ip.String() + "/32",
							Via:    ztNet.ExitNode,
						})
					}
				}

				for _, ip := range app.IPs {
					ztRoutes = append(ztRoutes, ZTRoute{
						Target: ip + "/32",
						Via:    ztNet.ExitNode,
					})
				}

				// app.Routes = ztRoutes
				ztNet.Apps[i].Routes = ztRoutes
			}

			if err := saveConfig(ztNet, configPath); err != nil {
				return err
			}

			ztNetworkConfig := ZTNetworkConfig{
				Routes: []ZTRoute{
					{Target: ztNet.LAN},
				},
			}

			for _, app := range ztNet.Apps {
				ztNetworkConfig.Routes = append(ztNetworkConfig.Routes, app.Routes...)
			}

			// sort.SliceStable(ztRoutes, func(i, j int) bool {
			// 	return ztRoutes[i].Target > ztRoutes[j].Target
			// })

			if debug {
				fmt.Print("DEBUG MODE: Skiping updating ZT controller")

			} else {
				ztToken := os.Getenv("ZT_TOKEN")
				if ztToken == "" {
					return fmt.Errorf("missing ZT_TOKEN environment variable")
				}

				if err := updateZTConfig(ZTConfig{Config: ztNetworkConfig}, ztNet.NetworkID, ztToken); err != nil {
					return err
				}
			}

			return nil
		},
	}

	rootCommand.PersistentFlags().StringVarP(&configPath, "config", "c", "./zt_routes.hjson", "config file")
	rootCommand.PersistentFlags().BoolVar(&debug, "debug", false, "option to update local config without updating ZT config")

	if err := rootCommand.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func loadConfig(path string) (ZTNetwork, error) {
	const operation = "loadConfig"

	configBytes, err := os.ReadFile(path)
	if err != nil {
		return ZTNetwork{}, fmt.Errorf("%s: %w", operation, err)
	}

	var config ZTNetwork
	err = hjson.Unmarshal(configBytes, &config)
	if err != nil {
		return ZTNetwork{}, fmt.Errorf("%s: %w", operation, err)
	}

	return config, err
}

func updateZTConfig(config ZTConfig, ztNetwork string, ztToken string) error {
	const operation = "updateZTConfig"

	bytesData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/network/%s", ZT_URL, ztNetwork), bytes.NewBuffer(bytesData))
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("token %s", ztToken))

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("%s:%w", "", err)

		}

		return fmt.Errorf("%s: %s", operation, string(body))
	}

	defer resp.Body.Close()

	return nil
}

func saveConfig(config ZTNetwork, configPath string) error {
	const operation = "saveConfig"

	bytesConfig, err := hjson.Marshal(config)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	defer file.Close()

	_, err = file.Write(bytesConfig)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	return nil
}
