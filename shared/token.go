package shared

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/uxland/gal-auth/model"
	"time"
)

var userCtxKey = "ctx-user-data"

var tokenKey = "ctx-auth-token"

var authorizationKey = "ctx-authorization"

func deserializeUserData(claims *jwt.MapClaims) (*model.TokenData, error) {
	var data map[string]interface{} = *claims
	user, b := data["user"]
	if !b {
		return nil, errors.New("user payload not found")
	}
	bytes, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	u := &model.TokenData{}
	err = json.Unmarshal(bytes, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

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
func setUserForContext(ctx context.Context, user *model.TokenData) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}
func SetClaimsForContext(ctx context.Context, claims *jwt.MapClaims) context.Context {
	ctx = context.WithValue(ctx, tokenKey, claims)
	userData, err := deserializeUserData(claims)
	if err != nil || userData == nil {
		return ctx
	}
	return setUserForContext(ctx, userData)
}

func SetAuthenticationToContext(ctx context.Context, authorizationType *AuthorizationType) context.Context {
	return context.WithValue(ctx, authorizationKey, authorizationType)
}

func GetAuthenticationFromContext(ctx context.Context) *AuthorizationType {
	auth, exists := ctx.Value(authorizationKey).(*AuthorizationType)
	if !exists {
		return nil
	}
	return auth
}

func GetContextData(ctx context.Context) map[string]interface{} {
	claims, exists := ctx.Value(tokenKey).(*jwt.MapClaims)
	if !exists || claims == nil {
		return nil
	}
	return *claims
}

func GetContextUser(ctx context.Context) *model.TokenData {
	user, exists := ctx.Value(userCtxKey).(*model.TokenData)
	if !exists || user == nil {
		return nil
	}
	return user
}

func TokenFactory(apiSecret string, duration time.Duration) func(user *model.TokenData) (string, error) {
	secret := []byte(apiSecret)
	return func(user *model.TokenData) (string, error) {
		claims := jwt.MapClaims{}
		claims["user"] = user
		if duration > time.Nanosecond*0 {
			claims["exp"] = time.Now().Add(duration).Unix()
		}
		claims["iat"] = time.Now().Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenStr, err := token.SignedString(secret)
		return tokenStr, err
	}
}
