package model

import (
	"cloud.google.com/go/datastore"
	"fmt"
)

func (u *UserData) Load(properties []datastore.Property) error {
	return datastore.LoadStruct(u, properties)
}

func (u *UserData) Save() ([]datastore.Property, error) {
	identities := make([]string, 0)
	for _, identity := range u.Identities {
		identities = append(identities, fmt.Sprintf("%s:%s", identity.ProviderID, identity.UserID))
	}
	u.XXXIdentities = identities
	return datastore.SaveStruct(u)
}
