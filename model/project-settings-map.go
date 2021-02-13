package model

import (
	"cloud.google.com/go/datastore"
)

type ProjectSettingsMap map[string]ProjectSettings

func (p *ProjectSettingsMap) Save() ([]datastore.Property, error) {
	props := make([]datastore.Property, 0)
	for key, settings := range *p {
		properties, err := datastore.SaveStruct(&settings)
		if err != nil {
			return nil, err
		}
		props = append(props, datastore.Property{Name: key, Value: &datastore.Entity{Properties: properties}})
	}
	return props, nil
}
func (p *ProjectSettingsMap) Load(props []datastore.Property) error {
	*p = ProjectSettingsMap{}
	for _, prop := range props {
		settings := &ProjectSettings{}
		e := prop.Value.(*datastore.Entity)
		err := datastore.LoadStruct(settings, e.Properties)
		if err != nil {
			return err
		}
		(*p)[prop.Name] = *settings
	}
	return nil
}
