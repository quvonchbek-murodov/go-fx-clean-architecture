package internal

import (
	"golang-project-structure/internal/dbconnections"
	"golang-project-structure/internal/repositories"
	"golang-project-structure/internal/services"

	"go.uber.org/fx"
)

var Module = fx.Options(
	dbconnections.Module,
	repositories.Module,
	services.Module,
)
