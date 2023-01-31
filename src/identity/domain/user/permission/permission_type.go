package permission

import "strings"

type Type int

const (
	Unknown   Type = 0
	ViewUsers Type = 1
	EditUsers Type = 2
	ViewRoles Type = 3
	EditRoles Type = 4
)

var permissionTypes = map[Type]string{
	Unknown:   "Unknown",
	ViewUsers: "View Users",
	EditUsers: "Edit Users",
	ViewRoles: "View Roles",
	EditRoles: "Edit Roles",
}

func (p Type) String() string {
	s, ok := permissionTypes[p]
	if !ok {
		return "Unknown"
	}
	return s
}

func ParsePermissionType(s string) (Type, bool) {
	for t, name := range permissionTypes {
		if strings.ToLower(s) == strings.ToLower(name) {
			return t, true
		}
	}
	return Unknown, false
}
