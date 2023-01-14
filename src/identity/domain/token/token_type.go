package token

import "strings"

type VerificationType int

const (
	Unknown VerificationType = iota
	Email
)

var verificationTypes = map[VerificationType]string{
	Unknown: "Unknown",
	Email:   "Email",
}

func (v VerificationType) String() string {
	s, ok := verificationTypes[v]
	if !ok {
		return "Unknown"
	}
	return s
}

func ParseVerificationType(s string) (VerificationType, bool) {
	for t, name := range verificationTypes {
		if strings.ToLower(s) == strings.ToLower(name) {
			return t, true
		}
	}

	return Unknown, false
}
