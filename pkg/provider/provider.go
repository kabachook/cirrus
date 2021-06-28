package provider

import (
	"inet.af/netaddr"
)

type Endpoint struct {
	Cloud string     `json:"cloud,omitempty"`
	IP    netaddr.IP `json:"ip,omitempty"`
	Type  string     `json:"type,omitempty"`
	Name  string     `json:"name,omitempty"`
}

type Provider interface {
	All() ([]Endpoint, error)
	Name() string
}
