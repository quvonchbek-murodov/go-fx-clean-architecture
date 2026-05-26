package handlers

import (
	userpb "golang-project-structure/genprotos/user"

	"go.uber.org/fx"
	"google.golang.org/grpc"
)

var Module = fx.Options(
	fx.Provide(NewUserHandler),

	fx.Invoke(RegisterUserHandler),
)

func RegisterUserHandler(server *grpc.Server, handler *UserHandler) {
	userpb.RegisterUserServiceServer(server, handler)
}
