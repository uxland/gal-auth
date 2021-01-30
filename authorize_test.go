package gal_auth

import (
	"context"
	"github.com/uxland/gal-auth/grpc"
	"github.com/uxland/gal-auth/model"
	"github.com/uxland/gal-auth/shared"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

func createUser() model.UserData {
	user := model.UserData{
		Profile: model.Profile{
			FullName:  "test user",
			PhotoLink: "http://cnd.test.jpg",
		},
		IsSuperUser: false,
		AccessTo: map[string]model.ProjectSettings{
			"p1": {
				Roles:       []string{"operator"},
				Date:        time.Now(),
				GrantedByID: "god",
			},
			"p2": {
				Roles:       []string{"admin"},
				Date:        time.Now(),
				GrantedByID: "god",
			},
		},
		Identities: map[string]model.Identity{
			"SAP": {UserID: "tt", ProviderID: "SAP"},
		},
	}
	return user
}

const (
	testApiSecret = "my-api-secret"
)

func setContext(user *model.UserData) (context.Context, error) {
	middleware := grpc.MiddlewareFactory(testApiSecret, shared.BearerSchema)
	tokenFactory := shared.TokenFactory(testApiSecret, time.Hour)

	payload, err := tokenFactory(user)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	ctx = grpc.Authenticate(ctx, shared.BearerSchema, payload)
	md, b := metadata.FromOutgoingContext(ctx)
	if !b {
		return nil, err
	}
	ctx = metadata.NewIncomingContext(ctx, md)
	ctx, err = middleware(ctx)
	if err != nil {
		return nil, err
	}
	return ctx, nil

}
func TestAssertRoleInProject(t *testing.T) {
	user := createUser()
	ctx, err := setContext(&user)
	err = AssertRoleInProject(ctx, "operator", "p1")
	if err != nil {
		t.Fail()
	}
	if !IsRoleInProject(ctx, "operator", "p1") {
		t.Fail()
	}
	err = AssertIsAdminInProject(ctx, "p2")
	if err != nil {
		t.Fail()
	}
	if !IsAdminInProject(ctx, "p2") {
		t.Fail()
	}
}
func TestAssertRoleForSuperUser(t *testing.T) {
	user := createUser()
	user.IsSuperUser = true
	ctx, err := setContext(&user)
	if err != nil {
		t.Error(err)
	}
	if err = AssertRoleInProject(ctx, "xxxxxxxxx", "yyyyyy"); err != nil {
		t.Fail()
	}
	if !IsRoleInProject(ctx, "xxxxxxxxx", "yyyyyy") {
		t.Fail()
	}
	if err = AssertIsAdminInProject(ctx, "xxxxxxxxx"); err != nil {
		t.Fail()
	}
	if !IsAdminInProject(ctx, "yyyyyy") {
		t.Fail()
	}
}
