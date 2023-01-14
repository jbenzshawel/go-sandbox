package token

import "github.com/google/uuid"

type Repository interface {
	GetToken(key uuid.UUID) *Token
	SaveToken(key uuid.UUID, token *Token)
	ClearToken(key uuid.UUID)
}
