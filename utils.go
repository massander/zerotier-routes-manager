package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

func loadConfig(path string) (Config, error) {
	const operation = "loadConfig"

	var emptyConfig Config

	configBytes, err := os.ReadFile(path)
	if err != nil {
		return emptyConfig, fmt.Errorf("%s: %w", operation, err)
	}

	var config []Group
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return emptyConfig, fmt.Errorf("%s: %w", operation, err)
	}

	return config, err
}

// syncRoutes push changes to ZT Network Contrller
func syncRoutes(config Config, network string, token string) error {
	const operation = "syncRoutes"

	ztRoutes := make([]ZTRoute, 0)

	// TODO: May be better find LAN from previous routes config?
	pools, err := getZTIPAssigmentPools(network, token)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	for _, pool := range pools {

		cidr, err := ipRangeToCIDR(pool.IPRangeStart, pool.IPRangeEnd)
		if err != nil {
			return fmt.Errorf("%s: %w", operation, err)
		}

		// fmt.Println("LAN:", cidr)

		ztRoutes = append(ztRoutes, ZTRoute{
			Target: cidr,
		})
	}

	// Default exitNode is defined under default group
	exitNode := config[0].ExitNode

	for _, group := range config {
		for _, route := range group.Routes {
			for _, address := range route.IPs {
				ztRoutes = append(ztRoutes,
					ZTRoute{
						Target: address + "/32",
						Via:    exitNode,
					})
			}
		}
	}

	// TODO: Sort is necessary???
	// sort.SliceStable(ztRoutes, func(i, j int) bool {
	// 	return ztRoutes[i].Target > ztRoutes[j].Target
	// })

	var payload = map[string]map[string]any{
		"config": {
			"routes": ztRoutes,
		},
	}

	payloadBytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	// fmt.Println(string(payloadBytes))

	URL := fmt.Sprintf("%s/network/%s", ZT_URL, network)

	request, err := http.NewRequest("POST", URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("%s:%w", operation, err)

		}

		return fmt.Errorf("%s: %s", operation, string(body))
	}

	return nil
}

func saveConfig(config Config, configPath string) error {
	const operation = "saveConfig"

	bytesConfig, err := json.MarshalIndent(config, "", "\t")
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

func ipRangeToCIDR(ipRangeStart, ipRangeEnd string) (string, error) {
	const operation = "ipRangeToCIDR"
	ipStart := net.ParseIP(ipRangeStart)
	ipEnd := net.ParseIP(ipRangeEnd)
	if ipStart == nil || ipEnd == nil {
		return "", fmt.Errorf("%s :invalid IP address", operation)
	}

	if ipStart.To4() != nil {
		ipStart = ipStart.To4()
	} else {
		ipStart = ipStart.To16()
	}

	if ipEnd.To4() != nil {
		ipEnd = ipEnd.To4()
	} else {
		ipEnd = ipEnd.To16()
	}

	mask := make([]byte, len(ipStart))

	for idx := range ipStart {
		mask[idx] = 255 - (ipStart[idx] ^ ipEnd[idx])
	}

	ipnet := net.IPNet{
		IP:   ipStart,
		Mask: mask,
	}

	ones, _ := ipnet.Mask.Size()

	cidr := fmt.Sprintf("%s/%d", ipStart.String(), uint8(ones))

	return cidr, nil
}

func getZTIPAssigmentPools(network string, token string) ([]ZTIPAssigmentsPool, error) {
	const operation = "getZTIPAssigmentPools"

	URL := fmt.Sprintf("%s/network/%s", ZT_URL, network)

	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return []ZTIPAssigmentsPool{}, fmt.Errorf("%s: %w", operation, err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("token %s", token))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return []ZTIPAssigmentsPool{}, fmt.Errorf("%s: %w", operation, err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []ZTIPAssigmentsPool{}, fmt.Errorf("%s:%w", operation, err)
	}

	if response.StatusCode != 200 {
		return []ZTIPAssigmentsPool{}, fmt.Errorf("%s: %s", operation, string(body))
	}

	var object map[string]json.RawMessage
	if err := json.Unmarshal(body, &object); err != nil {
		return []ZTIPAssigmentsPool{}, fmt.Errorf("%s: %s", operation, err)
	}

	var config map[string]json.RawMessage
	if err := json.Unmarshal(object["config"], &config); err != nil {
		return []ZTIPAssigmentsPool{}, fmt.Errorf("%s: %s", operation, err)
	}

	var ipAssigmentsPools []ZTIPAssigmentsPool
	if err := json.Unmarshal(config["ipAssignmentPools"], &ipAssigmentsPools); err != nil {
		return []ZTIPAssigmentsPool{}, fmt.Errorf("%s: %s", operation, err)
	}

	return ipAssigmentsPools, nil
}
