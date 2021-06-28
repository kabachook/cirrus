package server

import (
	"context"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/kabachook/cirrus/pkg/provider"
	"go.uber.org/zap"
)

type Server struct {
	ctx    context.Context
	logger *zap.Logger
	router *gin.Engine
	server *http.Server
}

type Config struct {
	Logger    *zap.Logger
	Providers []provider.Provider
	Server    *http.Server
}

func New(ctx context.Context, cfg Config) (*Server, error) {
	s := &Server{
		ctx:    ctx,
		logger: cfg.Logger,
		router: gin.Default(),
	}
	if err := s.init(cfg); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) init(cfg Config) error {
	s.router.Use(ginzap.Ginzap(s.logger, time.RFC3339, false))
	s.router.Use(ginzap.RecoveryWithZap(s.logger, true))

	api := s.router.Group("/v1")

	for _, provider := range cfg.Providers {
		GenerateProviderRoutes(api, provider)
		s.logger.Info("Provider added", zap.String("name", provider.Name()))
	}

	s.server = cfg.Server
	s.server.Handler = s.router

	return nil
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
