package model

import "testing"

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
