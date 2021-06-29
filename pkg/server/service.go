package server

import (
	"github.com/kabachook/cirrus/pkg/provider"
	"go.uber.org/zap"
)

func (s *Server) allEndpoints() ([]provider.Endpoint, error) {
	endpoints := make([]provider.Endpoint, 0)
	for _, provider := range s.providers {
		pEndpoints, err := provider.All()
		if err != nil {
			s.logger.Error("Failed to get endpoints", zap.Error(err))
			return nil, err
		}
		endpoints = append(endpoints, pEndpoints...)
	}
	return endpoints, nil
}
