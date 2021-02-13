package model

import (
	"fmt"
	"time"
)

const AdminRole = "admin"

type Identity struct {
	UserID      string
	ProviderID  string
	DisplayName string
}

type ProjectSettings struct {
	ID          string
	Roles       []string
	Properties  map[string]string
	GrantedByID string
	Date        time.Time
}

type Profile struct {
	FullName  string
	PhotoLink string
}

type UserData struct {
	ID          string
	Profile     Profile
	IsSuperUser bool
	AccessTo    map[string]ProjectSettings
	Identities  map[string]Identity
}

type TokenData struct {
	*UserData
	ThirdPartyTokens map[string]interface{}
}

func (u *UserData) String() string {
	return fmt.Sprintf("ID: %s; IsSuperdUser: %t; Identities: %+v", u.ID, u.IsSuperUser, u.Identities)
}
func (u *UserData) IsRoleInProject(projectID string, role string) bool {
	ps, found := u.AccessTo[projectID]
	if !found {
		return false
	}
	for _, r := range ps.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (u *UserData) IsAdminInProject(projectID string) bool {
	return u.IsRoleInProject(projectID, AdminRole)
}
