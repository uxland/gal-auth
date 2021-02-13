package model

import "cloud.google.com/go/datastore"

type IdentityMap map[string]Identity

func (m *IdentityMap) Save() ([]datastore.Property, error) {
	props := make([]datastore.Property, 0)
	for key, identity := range *m {
		idProps, err := datastore.SaveStruct(&identity)
		if err != nil {
			return nil, err
		}
		en := &datastore.Entity{Properties: idProps}
		props = append(props, datastore.Property{Name: key, Value: en})
	}
	return props, nil
}
func (m *IdentityMap) Load(props []datastore.Property) error {
	*m = IdentityMap{}
	for _, prop := range props {
		identity := &Identity{}
		e := prop.Value.(*datastore.Entity)
		err := datastore.LoadStruct(identity, e.Properties)
		if err != nil {
			return err
		}
		(*m)[prop.Name] = *identity
	}
	return nil
}
