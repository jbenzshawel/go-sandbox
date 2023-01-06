package storage

import (
	"github.com/google/uuid"
	"time"

	"github.com/jbenzshawel/go-sandbox/common/cache"
)

type VerificationTokenRepository struct {
	tokenCache *cache.ExpirationMap[uuid.UUID, string]
}

func NewVerificationTokenCache() *cache.ExpirationMap[uuid.UUID, string] {
	return cache.NewExpirationMap[uuid.UUID, string](20*time.Minute, false)
}

func NewVerificationTokenRepository(tokenCache *cache.ExpirationMap[uuid.UUID, string]) *VerificationTokenRepository {
	return &VerificationTokenRepository{
		tokenCache: tokenCache,
	}
}

func (r *VerificationTokenRepository) GetToken(userUUID uuid.UUID) string {
	return r.tokenCache.Get(userUUID)
}

func (r *VerificationTokenRepository) SaveToken(userUUID uuid.UUID, token string) {
	r.tokenCache.Set(userUUID, token)
}

func (r *VerificationTokenRepository) ClearToken(userUUID uuid.UUID) {
	r.tokenCache.Delete(userUUID)
}
