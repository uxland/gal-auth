package shared

import "context"

func AuthenticateInputToken(ctx context.Context, scheme AuthenticationScheme, token string, apiSecretBytes []byte) (context.Context, error) {
	ctx = SetAuthenticationToContext(ctx, &AuthorizationType{Scheme: scheme, Value: token})
	if scheme == BearerSchema {
		claims, err := VerifyToken(token, apiSecretBytes)
		if err != nil {
			return nil, err
		}
		ctx = SetClaimsForContext(ctx, claims)
	}
	return ctx, nil
}
