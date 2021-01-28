package gal_auth

import (
	"context"
	"github.com/uxland/ga-go-auth/grpc"
	"github.com/uxland/ga-go-auth/shared"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

type ProjectSettings struct {
	ID          string
	Roles       []string
	Properties  map[string]string
	GrantedByID string
	Date        time.Time
}
type Profile struct {
	FullName  string
	PhotoLink string
}

type UserDataDTO struct {
	Profile     Profile
	IsSuperUser bool
	AccessTo    map[string]ProjectSettings
	Identities  map[string]Identity
}
type Identity struct {
	UserID      string
	ProviderID  string
	DisplayName string
}

func createUser() UserDataDTO {
	user := UserDataDTO{
		Profile: Profile{
			FullName:  "test user",
			PhotoLink: "http://cnd.test.jpg",
		},
		IsSuperUser: false,
		AccessTo: map[string]ProjectSettings{
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
		Identities: map[string]Identity{
			"SAP": {UserID: "tt", ProviderID: "SAP"},
		},
	}
	return user
}

const (
	testApiSecret = "my-api-secret"
)

func setContext(user *UserDataDTO) (context.Context, error) {
	middleware := grpc.MiddlewareFactory(testApiSecret, shared.BearerSchema)
	tokenFactory := shared.TokenFactory(testApiSecret)

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
