package servers

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewGRPCServer),
	fx.Provide(NewGinEngine),
	fx.Provide(NewHTTPServer),

	fx.Invoke(StartGRPCServer),
	fx.Invoke(StartHTTPServer),
)
