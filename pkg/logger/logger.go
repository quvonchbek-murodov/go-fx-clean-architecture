package logger

import (
	"golang-project-structure/config"

	"go.uber.org/zap"
)

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	if cfg.AppDebug {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}
