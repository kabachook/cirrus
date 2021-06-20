package provider

import (
	"context"

	"inet.af/netaddr"
)

type Endpoint struct {
	IP   netaddr.IP
	Type string
	Name string
}

type Provider interface {
	New(context.Context) (*Provider, error)
	Configure(interface{}) error
	All() ([]*Endpoint, error)
}
