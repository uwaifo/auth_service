package iam

import (
	"strings"
)

type Permission struct {
	App string
	Scopes []string
	Permissions []string
}


// String() function will return the english name
// that we want out constant Day be recognized as
func (permission Permission) String() string {
	return "<Permission " + strings.Join(permission.Permissions, " ") + " />"
}

func (permission Permission) Raw() Permission  {
	return permission
}

func (permission Permission) HasPermission(userPermit string) bool {
	for _, permit := range permission.Permissions {
		if permit == userPermit {
			return true
		}
	}
	return false
}
