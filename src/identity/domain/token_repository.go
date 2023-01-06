package domain

import "github.com/google/uuid"

type TokenRepository interface {
	GetToken(key uuid.UUID) string
	SaveToken(key uuid.UUID, token string)
	ClearToken(key uuid.UUID)
}
