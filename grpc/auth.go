package grpc

import (
	"context"
	"fmt"
	"github.com/uxland/ga-go-auth/shared"
	"google.golang.org/grpc/metadata"
)

func Authenticate(ctx context.Context, scheme shared.AuthenticationScheme, auth string) context.Context {
	return metadata.AppendToOutgoingContext(ctx, shared.AuthorizationHeader, fmt.Sprintf("%s %s", scheme, auth))
}
