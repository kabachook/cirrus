package yc

import (
	"context"

	"github.com/kabachook/cirrus/pkg/provider"
	compute "github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	redis "github.com/yandex-cloud/go-genproto/yandex/cloud/mdb/redis/v1"
	vpc "github.com/yandex-cloud/go-genproto/yandex/cloud/vpc/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"go.uber.org/zap"
	"inet.af/netaddr"
)

const Name = "yc"

type Provider struct {
	ctx      context.Context
	sdk      *ycsdk.SDK
	logger   *zap.Logger
	folderId string
	zones    []string
}

type Config struct {
	Logger   *zap.Logger
	Token    string
	FolderID string
	Zones    []string
}

func New(ctx context.Context, cfg Config) (*Provider, error) {
	sdk, err := ycsdk.Build(ctx, ycsdk.Config{
		Credentials: ycsdk.OAuthToken(cfg.Token),
	})
	if err != nil {
		return nil, err
	}

	return &Provider{
		ctx:      ctx,
		sdk:      sdk,
		logger:   cfg.Logger,
		folderId: cfg.FolderID,
		zones:    cfg.Zones,
	}, nil
}

func (p *Provider) Name() string {
	return Name
}

func (p *Provider) Instances() ([]provider.Endpoint, error) {
	const typeName = "instance"
	var endpoints []provider.Endpoint

	res, err := p.sdk.Compute().Instance().List(p.ctx, &compute.ListInstancesRequest{
		FolderId: p.folderId,
	})
	if err != nil {
		return nil, err
	}

	p.logger.Debug("Response", zap.Any("instances", res.GetInstances()))

	for _, instance := range res.GetInstances() {
		for _, iface := range instance.GetNetworkInterfaces() {
			addr := iface.GetPrimaryV4Address()

			if ip := addr.GetAddress(); ip != "" {
				ip, err := netaddr.ParseIP(ip)
				if err != nil {
					return nil, err
				}
				endpoints = append(endpoints, provider.Endpoint{
					IP:   ip,
					Type: typeName,
					Name: instance.Name,
				})

			}

			if ip := addr.GetOneToOneNat().Address; ip != "" {
				ip, err := netaddr.ParseIP(ip)
				if err != nil {
					return nil, err
				}
				endpoints = append(endpoints, provider.Endpoint{
					IP:   ip,
					Type: typeName,
					Name: instance.Name,
				})
			}
		}
	}

	return endpoints, nil
}

func (p *Provider) Redis() ([]provider.Endpoint, error) {
	const typeName = "redis"
	var endpoints []provider.Endpoint

	// TODO: add pagination
	res, err := p.sdk.MDB().Redis().Cluster().List(p.ctx, &redis.ListClustersRequest{
		FolderId: p.folderId,
	})
	if err != nil {
		return nil, err
	}

	p.logger.Debug("Response", zap.Any("redis", res.GetClusters()))

	for _, cluster := range res.GetClusters() {
		hosts, err := p.sdk.MDB().Redis().Cluster().ListHosts(p.ctx, &redis.ListClusterHostsRequest{
			ClusterId: cluster.Id,
		})
		if err != nil {
			return nil, err
		}
		for _, host := range hosts.GetHosts() {
			endpoints = append(endpoints, provider.Endpoint{
				Type: typeName,
				Name: host.Name,
			})
		}
	}

	return endpoints, nil
}

func (p *Provider) Addresses() ([]provider.Endpoint, error) {
	const typeName = "address"
	var endpoints []provider.Endpoint

	// TODO: add pagination
	res, err := p.sdk.VPC().Address().List(p.ctx, &vpc.ListAddressesRequest{
		FolderId: p.folderId,
		Filter:   `type="EXTERNAL"`,
	})
	if err != nil {
		return nil, err
	}

	p.logger.Debug("Response", zap.Any("addresses", res.GetAddresses()))

	for _, address := range res.GetAddresses() {
		ip, err := netaddr.ParseIP(address.GetExternalIpv4Address().Address)
		if err != nil {
			return nil, err
		}
		var name string
		if address.GetName() != "" {
			name = address.Name
		} else {
			name = address.Id
		}
		endpoints = append(endpoints, provider.Endpoint{
			IP:   ip,
			Type: typeName,
			Name: name,
		})
	}

	return endpoints, nil
}

func (p *Provider) All() ([]provider.Endpoint, error) {
	globalFuncs := []func() ([]provider.Endpoint, error){
		p.Instances,
		p.Redis,
		p.Addresses,
	}
	endpoints := make([]provider.Endpoint, 0)

	p.logger.Debug("Getting endpoints", zap.String("folderId", p.folderId))

	// No zones for this call
	for _, f := range globalFuncs {
		resources, err := f()
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, resources...)
	}

	for i := range endpoints {
		endpoints[i].Cloud = Name
	}

	return endpoints, nil
}
