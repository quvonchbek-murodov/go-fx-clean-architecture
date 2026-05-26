package main

import (
	"golang-project-structure/app"
	"golang-project-structure/config"
	"golang-project-structure/internal"

	"go.uber.org/fx"
)

// @title          go-fx-clean-architecture API
// @version        1.0
// @description    HTTP API for the go-fx-clean-architecture service.
// @host           localhost:8080
// @BasePath       /api/v1
func main() {
	fx.New(
		config.Module,
		internal.Module,
		app.Module,

		fx.NopLogger,
	).Run()
}
