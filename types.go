package main

type Route struct {
	Domain string   `json:"domain"`
	IPs    []string `json:"ips"`
}

type Group struct {
	Name      string `json:"group"`
	IsEnabled bool   `json:"isEnabled,omitempty"`
	// Device (linux server) through which to route traffic.
	ExitNode string  `json:"exitNode,omitempty"`
	Routes   []Route `json:"routes"`
}

type Config []Group

type ZTRoute struct {
	Target string `json:"target"`
	Via    string `json:"via,omitempty"`
}

type ZTIPAssigmentsPool struct {
	IPRangeStart string `json:"ipRangeStart"`
	IPRangeEnd   string `json:"ipRangeEnd"`
}
