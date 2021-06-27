package provider

import (
	"context"

	"inet.af/netaddr"
)

type Endpoint struct {
	IP   netaddr.IP `json:"ip,omitempty"`
	Type string     `json:"type,omitempty"`
	Name string     `json:"name,omitempty"`
}

type Provider interface {
	New(context.Context) (*Provider, error)
	All() ([]*Endpoint, error)
}
