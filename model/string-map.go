package model

import "cloud.google.com/go/datastore"

type StringMap map[string]string

func (s *StringMap) Save() ([]datastore.Property, error) {
	props := make([]datastore.Property, 0)
	for key, value := range *s {
		props = append(props, datastore.Property{Name: key, Value: value})
	}
	return props, nil
}
func (s *StringMap) Load(props []datastore.Property) error {
	*s = StringMap{}
	for _, prop := range props {
		(*s)[prop.Name] = prop.Value.(string)
	}
	return nil
}
