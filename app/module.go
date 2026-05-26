package app

import (
	"golang-project-structure/app/controllers"
	"golang-project-structure/app/handlers"
	"golang-project-structure/app/servers"
	"golang-project-structure/pkg/logger"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(logger.NewLogger),
	handlers.Module,
	controllers.Module,
	servers.Module,
)
