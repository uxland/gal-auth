package model

import (
	"cloud.google.com/go/datastore"
	"reflect"
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	user := &UserData{
		ID:          "xxx",
		Profile:     Profile{},
		IsSuperUser: false,
		AccessTo:    nil,
		Identities: map[string]Identity{
			"SAP": {
				DisplayName: "hhhh",
				ProviderID:  "SAP",
				UserID:      "jr",
			}},
	}
	desc := user.String()
	if desc == "" {
		t.Fail()
	}
}
func TestUserSerialization(t *testing.T) {
	user := &UserData{
		ID: "my-user",
		Profile: Profile{
			FullName:  "my user name",
			PhotoLink: "http://my-picture.jpg",
		},
		IsSuperUser: false,
		AccessTo: ProjectSettingsMap{
			"pm": {
				ID:    "p,",
				Roles: []string{"admin", "avatar"},
				Properties: StringMap{
					"OperatorID": "ddddd",
					"Another":    "DDD",
				},
				GrantedByID: "adminsuperuser",
				Date:        time.Now(),
			},
			"so": {
				ID:    "so",
				Roles: []string{"admin", "boss"},
				Properties: StringMap{
					"ddd": "aaaa",
				},
				GrantedByID: "ffff",
				Date:        time.Now(),
			},
		},
		Identities: IdentityMap{
			"SAP": {
				UserID:      "sap-user",
				ProviderID:  "SAP",
				DisplayName: "sap user",
			},
			"AD": {
				UserID:      "ad-user",
				ProviderID:  "AD",
				DisplayName: "user in AD",
			},
		},
	}
	properties, err := user.Save()
	if err != nil {
		t.Fail()
	}
	u := &UserData{}
	err = datastore.LoadStruct(u, properties)
	if err != nil {
		t.Fail()
	}
	if reflect.DeepEqual(user, u) == false {
		t.Fail()
	}
}
