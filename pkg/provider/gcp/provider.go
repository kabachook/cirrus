package gcp

import (
	"context"

	"github.com/kabachook/cirrus/pkg/provider"
	"go.uber.org/zap"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"inet.af/netaddr"
)

type Provider struct {
	ctx     context.Context
	service *compute.Service
	logger  *zap.Logger
	project string
	zones   []string
}

type Config struct {
	Project string
	Options []option.ClientOption
	Logger  *zap.Logger
	Zones   []string
}

func New(ctx context.Context, cfg Config) (*Provider, error) {
	service, err := compute.NewService(ctx, cfg.Options...)
	if err != nil {
		return nil, err
	}

	return &Provider{
		ctx:     ctx,
		service: service,
		logger:  cfg.Logger,
		project: cfg.Project,
		zones:   cfg.Zones,
	}, nil
}

func (p *Provider) Name() string {
	return "gcp"
}

func (p *Provider) Instances(zone string) ([]provider.Endpoint, error) {
	var endpoints []provider.Endpoint

	req := p.service.Instances.List(p.project, zone)
	if err := req.Pages(p.ctx, func(page *compute.InstanceList) error {
		p.logger.Debug("Fetching endpoints", zap.String("zone", zone))
		for _, instance := range page.Items {
			for _, iface := range instance.NetworkInterfaces {
				ip, err := netaddr.ParseIP(iface.NetworkIP)
				if err != nil {
					return err
				}

				endpoints = append(endpoints, provider.Endpoint{
					IP:   ip,
					Name: instance.Name,
					Type: instance.Kind,
				})
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return endpoints, nil
}

func (p *Provider) All() ([]provider.Endpoint, error) {
	var endpoints []provider.Endpoint

	p.logger.Debug("Getting endpoints", zap.String("project", p.project))

	for _, zone := range p.zones {
		zoneInstances, err := p.Instances(zone)
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, zoneInstances...)
	}

	return endpoints, nil
}
