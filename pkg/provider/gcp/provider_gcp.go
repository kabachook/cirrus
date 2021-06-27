package gcp

import (
	"context"

	"github.com/kabachook/cirrus/pkg/config"
	"github.com/kabachook/cirrus/pkg/provider"
	"go.uber.org/zap"
	"google.golang.org/api/compute/v1"
	"inet.af/netaddr"
)

type ProviderGCP struct {
	ctx     context.Context
	service *compute.Service
	logger  *zap.Logger
	project string
}

func New(ctx context.Context, cfg config.ConfigGCP) (*ProviderGCP, error) {
	service, err := compute.NewService(ctx, cfg.Options...)
	if err != nil {
		return nil, err
	}

	return &ProviderGCP{
		ctx:     ctx,
		service: service,
		logger:  cfg.Logger,
		project: cfg.Project,
	}, nil
}

func (p *ProviderGCP) Instances(zone string) ([]provider.Endpoint, error) {
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

func (p *ProviderGCP) All() ([]provider.Endpoint, error) {
	var endpoints []provider.Endpoint

	p.logger.Debug("Getting endpoints", zap.String("project", p.project))

	var zones []string
	req := p.service.Zones.List(p.project)
	if err := req.Pages(p.ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			zones = append(zones, zone.Description)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	for _, zone := range zones {
		zoneInstances, err := p.Instances(zone)
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, zoneInstances...)
	}

	return endpoints, nil
}
