package config

import (
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type ConfigGCP struct {
	Project string
	Options []option.ClientOption
	Logger  *zap.Logger
}
