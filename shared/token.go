package shared

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var userCtxKey = "user-auth"

func VerifyToken(token string, apiSecret []byte) (*jwt.MapClaims, error) {
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
func SetClaimsForContext(ctx context.Context, claims *jwt.MapClaims) context.Context {
	return context.WithValue(ctx, userCtxKey, claims)
}
func GetContextData(ctx context.Context) map[string]interface{} {
	claims, exists := ctx.Value(userCtxKey).(*jwt.MapClaims)
	if !exists || claims == nil {
		return nil
	}
	return *claims
}

func TokenFactory(apiSecret string) func(payload interface{}) (string, error) {
	secret := []byte(apiSecret)
	return func(payload interface{}) (string, error) {
		claims := jwt.MapClaims{}
		claims["payload"] = payload
		claims["exp"] = time.Now().Add(time.Hour * 8).Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString(secret)
		return tokenStr, err
	}
}
