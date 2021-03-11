package model

import (
	"testing"
	"time"
)

func TestName(t *testing.T) {
	dt := &UserData{
		ID: "fff",
		Profile: Profile{
			FullName:  "fddf",
			PhotoLink: "dsdsd",
		},
		IsSuperUser: false,
		AccessTo: map[string]ProjectSettings{
			"pm": {
				ID:          "deee",
				Roles:       []string{"ewewewe"},
				GrantedByID: "me",
				Properties:  nil,
				Date:        time.Now(),
			},
		},
		Identities: map[string]Identity{"SAP": {
			UserID:      "dd",
			ProviderID:  "ss",
			DisplayName: "ds",
		}},
		XXXIdentities: []string{"SAP:SSS"},
	}
	save, err := dt.Save()
	if err != nil {
		t.Fail()
	}
	if save == nil {
		t.Fail()
	}
}
