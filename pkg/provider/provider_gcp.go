package provider

import (
	"context"
	"net/http"

	"github.com/kabachook/cirrus/pkg/config"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"inet.af/netaddr"
)

type ProviderGCP struct {
	ctx     context.Context
	client  *http.Client
	project string
}

func (p *ProviderGCP) New(ctx context.Context) (*ProviderGCP, error) {
	return &ProviderGCP{
		ctx: ctx,
	}, nil
}

func (p *ProviderGCP) Configure(cfg config.ConfigGCP) error {
	p.project = cfg.Project
	return nil
}

func (p *ProviderGCP) All() ([]Endpoint, error) {
	var endpoints []Endpoint

	service, err := compute.NewService(p.ctx, option.WithHTTPClient(p.client))
	if err != nil {
		return nil, err
	}

	var zones []string
	req := service.Zones.List(p.project)
	if err := req.Pages(p.ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			zones = append(zones, zone.Description)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	for _, zone := range zones {
		req := service.Instances.List(p.project, zone)
		if err := req.Pages(p.ctx, func(page *compute.InstanceList) error {
			for _, instance := range page.Items {
				for _, iface := range instance.NetworkInterfaces {
					ip, err := netaddr.ParseIP(iface.NetworkIP)
					if err != nil {
						return err
					}

					endpoints = append(endpoints, Endpoint{
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
	}

	return endpoints, nil
}
