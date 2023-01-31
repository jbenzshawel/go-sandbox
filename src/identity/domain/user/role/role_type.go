package role

import (
	"strings"
)

type Type int

const (
	Unknown Type = 0
	Admin   Type = 1
)

var roleTypes = map[Type]string{
	Unknown: "Unknown",
	Admin:   "Admin",
}

func (r Type) String() string {
	s, ok := roleTypes[r]
	if !ok {
		return "Unknown"
	}
	return s
}

func ParsePermissionType(s string) (Type, bool) {
	for t, name := range roleTypes {
		if strings.ToLower(s) == strings.ToLower(name) {
			return t, true
		}
	}
	return Unknown, false
}
