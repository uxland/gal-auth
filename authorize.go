package gal_auth

import (
	"context"
	"github.com/stoewer/go-strcase"
	"github.com/uxland/ga-go-auth/shared"
	"reflect"
)

const AdminRole = "admin"

type caseConverter = func(string) string

var caseConverters = []caseConverter{strcase.UpperCamelCase, strcase.LowerCamelCase, strcase.SnakeCase, strcase.UpperSnakeCase, strcase.KebabCase, strcase.UpperKebabCase}

func getFieldInCase(data map[string]interface{}, name string, converter func(s string) string) (interface{}, bool) {
	name = converter(name)
	val, b := data[name]
	return val, b
}

func getField(data map[string]interface{}, fieldName string) (interface{}, bool) {
	var field interface{}
	var b bool
	field, b = data[fieldName]
	if b {
		return field, b
	}
	for _, converter := range caseConverters {
		field, b = getFieldInCase(data, fieldName, converter)
		if b {
			return field, b
		}
	}
	return nil, false
}

func isSuperUser(user map[string]interface{}) bool {
	value, exists := user["IsSuperUser"]
	if !exists {
		return false
	}
	if tp := reflect.TypeOf(value); tp.AssignableTo(reflect.TypeOf(true)) {
		return value.(bool)
	}
	return false
}

func IsRoleInProject(ctx context.Context, role, project string) bool {
	data := shared.GetContextData(ctx)
	if data == nil {
		return false
	}
	p, exists := getField(data, "payload")
	if !exists {
		return false
	}

	user := p.(map[string]interface{})
	if user == nil {
		return false
	}
	if isSuperUser(user) {
		return true
	}

	access, exists := getField(user, "AccessTo")
	if !exists {
		return false
	}
	accessTo := access.(map[string]interface{})

	prj, exists := accessTo[project].(map[string]interface{})
	if !exists {
		return false
	}
	rolesField, exists := getField(prj, "Roles")
	if !exists {
		return false
	}
	if rolesType := reflect.TypeOf(rolesField); rolesType.AssignableTo(reflect.TypeOf([]interface{}{})) {
		roles := rolesField.([]interface{})
		for _, r := range roles {
			if r == role {
				return true
			}
		}
	}
	return false

}
func IsAdminInProject(ctx context.Context, project string) bool {
	return IsRoleInProject(ctx, AdminRole, project)
}

func AssertRoleInProject(ctx context.Context, role, project string) error {
	if isRole := IsRoleInProject(ctx, role, project); !isRole {
		return shared.UnauthorizedError{}
	}
	return nil
}
func AssertIsAdminInProject(ctx context.Context, project string) error {
	return AssertRoleInProject(ctx, AdminRole, project)
}
