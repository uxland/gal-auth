package grpc_auth

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
)

func AuthenticateCtx(ctx context.Context, scheme AuthenticationScheme, auth string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", fmt.Sprintf("%s %s", scheme, auth))
}
