package gal_auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/stoewer/go-strcase"
	"github.com/uxland/gal-auth/model"
	"github.com/uxland/gal-auth/shared"
	"reflect"
)

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

func getUserData(ctx context.Context) (*model.UserData, map[string]interface{}, bool) {
	u := shared.GetContextUser(ctx)
	if u != nil {
		return u, nil, true
	}
	data := shared.GetContextData(ctx)
	if data == nil {
		return nil, nil, false
	}
	p, exists := getField(data, "user")
	if !exists {
		return nil, nil, false
	}
	user := p.(map[string]interface{})
	if user == nil {
		return nil, nil, false
	}
	return nil, user, true
}

func isSuperUser(user *model.UserData, mapUser map[string]interface{}) bool {
	if user != nil {
		return user.IsSuperUser
	}
	value, b := getField(mapUser, "IsSuperUser")
	if !b {
		return false
	}
	if tp := reflect.TypeOf(value); tp.AssignableTo(reflect.TypeOf(true)) {
		return value.(bool)
	}
	return false
}

func IsRoleInProject(ctx context.Context, role, project string) bool {
	user, mapUser, exists := getUserData(ctx)
	if !exists {
		return false
	}

	if isSuperUser(user, mapUser) {
		return true
	}
	if user != nil {
		return user.IsRoleInProject(project, role)
	}
	access, exists := getField(mapUser, "AccessTo")
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
	return IsRoleInProject(ctx, model.AdminRole, project)
}

func AssertRoleInProject(ctx context.Context, role, project string) error {
	if isRole := IsRoleInProject(ctx, role, project); !isRole {
		return shared.UnauthorizedError{}
	}
	return nil
}
func AssertIsAdminInProject(ctx context.Context, project string) error {
	return AssertRoleInProject(ctx, model.AdminRole, project)
}
func GetIncomingUserID(ctx context.Context) string {
	user, mapUser, found := getUserData(ctx)
	if !found {
		return ""
	}
	if user != nil {
		return user.ID
	}
	id, b := getField(mapUser, "ID")
	if !b {
		return ""
	}
	return id.(string)

}

func GetIncomingUserDescription(ctx context.Context) string {
	u, _, b := getUserData(ctx)
	if !b {
		return ""
	}
	if u != nil {
		return u.String()
	}
	return fmt.Sprintf("ID: %s", GetIncomingUserID(ctx))
}

func IsSuperUser(ctx context.Context) bool {
	user, mapUser, b := getUserData(ctx)
	if !b {
		return false
	}
	return isSuperUser(user, mapUser)
}

func AssertIsSuperUser(ctx context.Context) error {
	if !IsSuperUser(ctx) {
		return errors.New("user is not super user")
	}
	return nil
}

type Identity struct {
	UserID      string
	ProviderID  string
	DisplayName string
}
type UserInfo struct {
	ID          string
	IsSuperUser string
}
