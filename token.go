package grpc_auth

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
)
var userCtxKey = "user-auth"

func verifyToken(token string, apiSecret []byte) (*jwt.MapClaims, error) {
	var verifiedToken, tokenError = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return apiSecret, nil
	})
	if tokenError != nil {
		return nil, errors.New("error verifying token")
	}
	claims := verifiedToken.Claims.(jwt.MapClaims)
	tokenError = claims.Valid()
	if tokenError != nil {
		return nil, tokenError
	}
	return &claims, tokenError
}
func setClaimsForContext(ctx context.Context, claims *jwt.MapClaims) context.Context {
	return context.WithValue(ctx, userCtxKey, claims)
}
func GetContextData(ctx context.Context) map[string] interface{}{
	claims := ctx.Value(userCtxKey).(*jwt.MapClaims)
	return *claims
}