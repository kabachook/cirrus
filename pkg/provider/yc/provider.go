package yc

import (
	"context"

	"github.com/kabachook/cirrus/pkg/provider"
	compute "github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"go.uber.org/zap"
	"inet.af/netaddr"
)

const name = "yc"

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
	return name
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

	p.logger.Debug("Response", zap.Any("instances", res.Instances))

	for _, instance := range res.Instances {
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

func (p *Provider) All() ([]provider.Endpoint, error) {
	endpoints := make([]provider.Endpoint, 0)

	p.logger.Debug("Getting endpoints", zap.String("folderId", p.folderId))

	// No zones for this call
	// for _, zone := range p.zones {
	zoneInstances, err := p.Instances()
	if err != nil {
		return nil, err
	}
	endpoints = append(endpoints, zoneInstances...)
	// }

	for i := range endpoints {
		endpoints[i].Cloud = name
	}

	return endpoints, nil
}
