package database

import (
	"github.com/kabachook/cirrus/pkg/provider"
)

type Snapshot struct {
	Timestamp int64               `json:"timestamp,omitempty"`
	Endpoints []provider.Endpoint `json:"endpoints,omitempty"`
}

type Database interface {
	Open() error
	Close() error
	Store(int64, []provider.Endpoint) error
	List() ([]Snapshot, error)
	Get(int64) (*Snapshot, error)
}
