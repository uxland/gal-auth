package shared

import (
	"context"
	. "github.com/onsi/gomega"
	"github.com/uxland/gal-auth/model"
	"testing"
	"time"
)

func TestTokenFactory(t *testing.T) {
	RegisterTestingT(t)
	factory := TokenFactory("my-secret", time.Hour*8)
	user := &model.UserData{
		ID: "xxx",
		Profile: model.Profile{
			FullName:  "jr jr jr",
			PhotoLink: "http://imgs/dallas/jr.jpg",
		},
		IsSuperUser: false,
		AccessTo: map[string]model.ProjectSettings{
			"p1": {
				ID:          "p1",
				Roles:       []string{"admin"},
				Properties:  map[string]string{"operator": "125"},
				GrantedByID: "js",
				Date:        time.Now(),
			},
		},
		Identities: map[string]model.Identity{
			"SAP": {
				DisplayName: "hhhh",
				ProviderID:  "SAP",
				UserID:      "jr",
			},
			"AD": {
				DisplayName: "jr in AD",
				ProviderID:  "AD",
				UserID:      "adjr",
			},
		},
	}
	token, err := factory(user)
	if err != nil {
		t.Fail()
	}
	claims, err := VerifyToken(token, []byte("my-secret"))
	if err != nil {
		t.Fail()
	}
	ctx := SetClaimsForContext(context.Background(), claims)
	contextUser := GetContextUser(ctx)
	if contextUser == nil {
		t.Fail()
	}
	res := Expect(user).To(Equal(user))
	if !res {
		t.Fail()
	}

}
