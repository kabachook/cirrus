package server

import (
	"context"
	"net/http"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/kabachook/cirrus/pkg/database"
	"github.com/kabachook/cirrus/pkg/provider"
	"go.uber.org/zap"
)

type Server struct {
	ctx        context.Context
	logger     *zap.Logger
	router     *gin.Engine
	db         database.Database
	server     *http.Server
	providers  map[string]provider.Provider
	ScanPeriod time.Duration
}

type Config struct {
	Logger     *zap.Logger
	Providers  []provider.Provider
	Database   database.Database
	Server     *http.Server
	ScanPeriod time.Duration
}

func New(ctx context.Context, cfg Config) (*Server, error) {
	s := &Server{
		ctx:        ctx,
		logger:     cfg.Logger,
		router:     gin.Default(),
		db:         cfg.Database,
		providers:  make(map[string]provider.Provider),
		ScanPeriod: cfg.ScanPeriod,
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

	// /v1/<name>/all routes
	s.providersRoutes(api, cfg.Providers)

	// Premature optimization?
	providerNames := make([]string, 0, len(cfg.Providers))
	for _, provider := range cfg.Providers {
		s.logger.Debug("Adding provider to map", zap.String("name", provider.Name()))
		s.providers[provider.Name()] = provider
		providerNames = append(providerNames, provider.Name())
	}

	api.GET("/available", func(c *gin.Context) {
		c.JSON(http.StatusOK, providerNames)
	})

	api.GET("/all", func(c *gin.Context) {
		endpoints, err := s.allEndpoints()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, endpoints)
	})

	api.GET("/snapshots", func(c *gin.Context) {
		snapshots, err := s.db.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, snapshots)
	})

	api.POST("/snapshot/new", func(c *gin.Context) {
		endpoints, err := s.allEndpoints()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = s.db.Store(time.Now().Unix(), endpoints)
		if err != nil {
			s.logger.Error("Failed to save snapshot", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.Status(200)
	})

	s.server = cfg.Server
	s.server.Handler = s.router

	return nil
}

func (s *Server) Run() error {
	go s.runScanner()
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
