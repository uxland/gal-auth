package grpc

import (
	"context"
	grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/uxland/gal-auth/shared"
)

func MiddlewareFactory(apiSecret string, schemes ...shared.AuthenticationScheme) grpcAuth.AuthFunc {
	apiSecretBytes := []byte(apiSecret)
	if len(schemes) == 0 {
		schemes = []shared.AuthenticationScheme{shared.BearerSchema, shared.BasicSchema}
	}
	return func(ctx context.Context) (context.Context, error) {

		for _, scheme := range schemes {
			token, err := grpcAuth.AuthFromMD(ctx, scheme)
			if err != nil || token == "" {
				continue
			}
			ctx = shared.SetAuthenticationToContext(ctx, &shared.AuthorizationType{Scheme: scheme, Value: token})
			if scheme == shared.BearerSchema {
				claims, err := shared.VerifyToken(token, apiSecretBytes)
				if err != nil {
					return nil, err
				}
				ctx = shared.SetClaimsForContext(ctx, claims)
			}
			return ctx, nil
		}
		return ctx, nil
	}
}
