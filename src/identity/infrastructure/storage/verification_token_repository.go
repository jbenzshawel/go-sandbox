package storage

import (
	"time"

	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/common/cache"
	"github.com/jbenzshawel/go-sandbox/identity/domain/token"
)

type VerificationTokenRepository struct {
	tokenCache *cache.ExpirationMap[uuid.UUID, *token.Token]
}

func NewVerificationTokenCache() *cache.ExpirationMap[uuid.UUID, *token.Token] {
	return cache.NewExpirationMap[uuid.UUID, *token.Token](20*time.Minute, false)
}

func NewVerificationTokenRepository(tokenCache *cache.ExpirationMap[uuid.UUID, *token.Token]) *VerificationTokenRepository {
	return &VerificationTokenRepository{
		tokenCache: tokenCache,
	}
}

func (r *VerificationTokenRepository) GetToken(userUUID uuid.UUID) *token.Token {
	return r.tokenCache.Get(userUUID)
}

func (r *VerificationTokenRepository) SaveToken(userUUID uuid.UUID, token *token.Token) {
	r.tokenCache.Set(userUUID, token)
}

func (r *VerificationTokenRepository) ClearToken(userUUID uuid.UUID) {
	r.tokenCache.Delete(userUUID)
}
