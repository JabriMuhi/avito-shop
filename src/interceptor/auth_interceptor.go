package interceptor

import (
	"avito-shop/src/auth"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func authInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if info.FullMethod == "/avito.AvitoShop/Authenticate" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		values := md.Get("authorization")
		if len(values) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		token := values[0]
		claims, err := auth.ValidateJWT(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token")
		}

		ctx = context.WithValue(ctx, "user", claims.UserID)
		return handler(ctx, req)
	}
}

func CreateAuthInterceptor() grpc.UnaryServerInterceptor {
	return authInterceptor()
}
