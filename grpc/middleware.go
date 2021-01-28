package grpc

import (
	"context"
	grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/uxland/ga-go-auth/shared"
)

func MiddlewareFactory(apiSecret, scheme shared.AuthenticationScheme) grpcAuth.AuthFunc {
	apiSecretBytes := []byte(apiSecret)
	return func(ctx context.Context) (context.Context, error) {
		token, err := grpcAuth.AuthFromMD(ctx, scheme)
		if err != nil || token == "" {
			return ctx, nil
		}
		claims, err := shared.VerifyToken(token, apiSecretBytes)
		if err != nil {
			return nil, err
		}
		ctx = shared.SetClaimsForContext(ctx, claims)
		return ctx, nil
	}
}
