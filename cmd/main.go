package main

import (
	"golang-project-structure/app"
	"golang-project-structure/config"
	"golang-project-structure/internal"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		internal.Module,
		app.Module,

		fx.NopLogger,
	).Run()
}
