package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kabachook/cirrus/pkg/provider"
)

// GenerateProviderRoutes generates routes for provided provider with prefix `p.Name()`
func GenerateProviderRoutes(r *gin.RouterGroup, p provider.Provider) {
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
