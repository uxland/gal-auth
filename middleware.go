package grpc_auth

import (
	"context"
	auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

func MiddlewareFactory(apiSecret, scheme AuthenticationScheme) auth.AuthFunc {
	apiSecretBytes := []byte(apiSecret)
	return func(ctx context.Context) (context.Context, error) {
		token, err := auth.AuthFromMD(ctx, scheme)
		if err != nil || token == ""{
			return ctx, nil
		}
		claims, err := verifyToken(token, apiSecretBytes)
		if err != nil{
			return nil, err
		}
		ctx = setClaimsForContext(ctx, claims)
		return ctx, nil
	}
}