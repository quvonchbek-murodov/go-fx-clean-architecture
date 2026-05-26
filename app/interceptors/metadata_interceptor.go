package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func MetadataInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			for k, v := range md {
				if len(v) > 0 {
					ctx = context.WithValue(ctx, k, v[0])
				}
			}
		}
		return handler(ctx, req)
	}
}
