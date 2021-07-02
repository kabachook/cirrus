package config

import (
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	env := os.Getenv("CIRRUS_ENV")
	if env != "production" {
		Logger, _ = zap.NewDevelopment()
	} else {
		gin.SetMode(gin.ReleaseMode)
		Logger, _ = zap.NewProduction()
	}
}
