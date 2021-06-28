package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kabachook/cirrus/pkg/provider"
	"go.uber.org/zap"
)

// generateProviderRoutes generates routes for provided provider with prefix `p.Name()`
func generateProviderRoutes(r *gin.RouterGroup, p provider.Provider) {
	group := r.Group(fmt.Sprintf("/%s", p.Name()))
	group.GET("/all", func(c *gin.Context) {
		endpoints, err := p.All()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, endpoints)
	})
}

func (s *Server) providersRoutes(r *gin.RouterGroup, providers []provider.Provider) {
	for _, provider := range providers {
		generateProviderRoutes(r, provider)
		s.logger.Info("Provider added", zap.String("name", provider.Name()))
	}
}
