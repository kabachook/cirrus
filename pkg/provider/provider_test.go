package provider_test

import (
	"encoding/json"
	"testing"

	"github.com/kabachook/cirrus/pkg/provider"
	"github.com/stretchr/testify/assert"
	"inet.af/netaddr"
)

func TestMarshall(t *testing.T) {
	ip, err := netaddr.ParseIP("127.0.0.1")
	if err != nil {
		t.Error(err)
	}

	endpoint := &provider.Endpoint{
		IP:   ip,
		Type: "instance",
		Name: "test-instance",
	}

	b, err := json.Marshal(endpoint)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, string(b), `{"ip":"127.0.0.1","type":"instance","name":"test-instance"}`, "Endpoint.Marshall failed")
}
