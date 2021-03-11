package websockets

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/google/martian/log"
	"github.com/uxland/gal-auth/shared"
	"strings"
)

func MiddlewareFactory(apiSecret string) transport.WebsocketInitFunc {
	apiSecretBytes := []byte(apiSecret)
	return func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
		authorization := initPayload.Authorization()
		if authorization == "" {
			log.Errorf("user not authenticated \n")
			return nil, errors.New("access denied")
		}
		tokenSplit := strings.Split(authorization, " ")

		scheme := tokenSplit[0]
		auth := tokenSplit[1]
		ctx, err := shared.AuthenticateInputToken(ctx, scheme, auth, apiSecretBytes)
		if err != nil {
			log.Errorf("authentication error %s \n", err.Error())
			return nil, errors.New("access denied")
		}
		return ctx, nil
	}
}
