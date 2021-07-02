package gcp

import (
	"context"

	"github.com/kabachook/cirrus/pkg/provider"
	"go.uber.org/zap"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"inet.af/netaddr"
)

const Name = "gcp"

type Provider struct {
	ctx        context.Context
	service    *compute.Service
	logger     *zap.Logger
	project    string
	zones      []string
	aggregated bool
}

type Config struct {
	Project    string
	Options    []option.ClientOption
	Logger     *zap.Logger
	Zones      []string
	Aggregated bool
}

func New(ctx context.Context, cfg Config) (*Provider, error) {
	service, err := compute.NewService(ctx, cfg.Options...)
	if err != nil {
		return nil, err
	}

	return &Provider{
		ctx:        ctx,
		service:    service,
		logger:     cfg.Logger,
		project:    cfg.Project,
		zones:      cfg.Zones,
		aggregated: cfg.Aggregated,
	}, nil
}

func (p *Provider) Name() string {
	return Name
}

func processInstanceList(instances []*compute.Instance) ([]provider.Endpoint, error) {
	var endpoints []provider.Endpoint

	for _, instance := range instances {
		for _, iface := range instance.NetworkInterfaces {
			ip, err := netaddr.ParseIP(iface.NetworkIP)
			if err != nil {
				return nil, err
			}

			endpoints = append(endpoints, provider.Endpoint{
				IP:   ip,
				Name: instance.Name,
				Type: instance.Kind,
			})

			for _, aconf := range iface.AccessConfigs {
				if aconf.NatIP != "" {
					ip, err := netaddr.ParseIP(aconf.NatIP)
					if err != nil {
						return nil, err
					}
					endpoints = append(endpoints, provider.Endpoint{
						IP:   ip,
						Name: instance.Name,
						Type: instance.Kind,
					})
				}
			}
		}
	}

	return endpoints, nil
}

func (p *Provider) Instances(zone string) ([]provider.Endpoint, error) {
	var endpoints []provider.Endpoint

	req := p.service.Instances.List(p.project, zone)
	if err := req.Pages(p.ctx, func(page *compute.InstanceList) error {
		p.logger.Debug("Response", zap.Any("instances", page.Items))
		p.logger.Debug("Fetching endpoints", zap.String("zone", zone))
		pageEndpoints, err := processInstanceList(page.Items)
		if err != nil {
			return err
		}
		endpoints = append(endpoints, pageEndpoints...)
		return nil
	}); err != nil {
		return nil, err
	}

	return endpoints, nil
}

func (p *Provider) InstancesAggregated() ([]provider.Endpoint, error) {
	var endpoints []provider.Endpoint

	req := p.service.Instances.AggregatedList(p.project)
	if err := req.Pages(p.ctx, func(page *compute.InstanceAggregatedList) error {
		p.logger.Debug("Response", zap.Any("instances", page.Items))

		for _, scoped := range page.Items {
			scopedEndpoints, err := processInstanceList(scoped.Instances)
			if err != nil {
				return err
			}
			endpoints = append(endpoints, scopedEndpoints...)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return endpoints, nil
}

func processAddressList(addresses []*compute.Address) ([]provider.Endpoint, error) {
	var endpoints []provider.Endpoint

	for _, address := range addresses {
		ip, err := netaddr.ParseIP(address.Address)
		if err != nil {
			return nil, err
		}

		endpoints = append(endpoints, provider.Endpoint{
			IP:   ip,
			Name: address.Name,
			Type: address.Kind,
		})
	}

	return endpoints, nil
}

func (p *Provider) AddressesAggregated() ([]provider.Endpoint, error) {
	var endpoints []provider.Endpoint

	req := p.service.Addresses.AggregatedList(p.project)
	if err := req.Pages(p.ctx, func(page *compute.AddressAggregatedList) error {
		p.logger.Debug("Response", zap.Any("addresses", page.Items))

		for _, scoped := range page.Items {
			scopedEndpoints, err := processAddressList(scoped.Addresses)
			if err != nil {
				return err
			}
			endpoints = append(endpoints, scopedEndpoints...)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (p *Provider) GlobalAddresses() ([]provider.Endpoint, error) {
	var endpoints []provider.Endpoint

	req := p.service.GlobalAddresses.List(p.project)
	if err := req.Pages(p.ctx, func(page *compute.AddressList) error {
		p.logger.Debug("Response", zap.Any("addresses", page.Items))

		pageEndpoints, err := processAddressList(page.Items)
		if err != nil {
			return err
		}
		endpoints = append(endpoints, pageEndpoints...)
		return nil
	}); err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (p *Provider) All() ([]provider.Endpoint, error) {
	resourcesFuncs := []func(string) ([]provider.Endpoint, error){
		p.Instances,
		// TODO: p.Addresses
	}
	aggregatedResourcesFuncs := []func() ([]provider.Endpoint, error){
		p.InstancesAggregated,
		p.AddressesAggregated,
	}
	globalFuncs := []func() ([]provider.Endpoint, error){
		p.GlobalAddresses,
	}

	endpoints := make([]provider.Endpoint, 0)

	p.logger.Debug("Getting endpoints", zap.String("project", p.project))

	if p.aggregated {
		for _, getResource := range aggregatedResourcesFuncs {
			resourceEndpoints, err := getResource()
			if err != nil {
				return nil, err
			}
			endpoints = append(endpoints, resourceEndpoints...)
		}
	} else {
		for _, getResource := range resourcesFuncs {
			for _, zone := range p.zones {
				zoneInstances, err := getResource(zone)
				if err != nil {
					return nil, err
				}
				endpoints = append(endpoints, zoneInstances...)
			}
		}
	}

	for _, getResource := range globalFuncs {
		zoneInstances, err := getResource()
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, zoneInstances...)
	}

	for i := range endpoints {
		endpoints[i].Cloud = Name
	}

	return endpoints, nil
}
