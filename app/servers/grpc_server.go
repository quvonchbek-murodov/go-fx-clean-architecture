package servers

import (
	"context"
	"fmt"
	"net"

	"golang-project-structure/app/interceptors"
	"golang-project-structure/config"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer() *grpc.Server {
	return grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.MetadataInterceptor(),
		),
	)
}

func StartGRPCServer(lc fx.Lifecycle, server *grpc.Server, cfg *config.Config, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			reflection.Register(server)

			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
			if err != nil {
				return fmt.Errorf("grpc listen: %w", err)
			}

			go func() {
				log.Info("grpc server listening", zap.String("addr", lis.Addr().String()))
				if err := server.Serve(lis); err != nil {
					log.Error("grpc serve error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("grpc server stopping")
			server.GracefulStop()
			return nil
		},
	})
}
