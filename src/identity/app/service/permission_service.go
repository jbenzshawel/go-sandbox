package service

import (
	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/identity/domain/user"
	"github.com/jbenzshawel/go-sandbox/identity/domain/user/permission"
)

type PermissionService struct {
	userRepo user.Repository
}

func NewPermissionService(userRepo user.Repository) *PermissionService {
	return &PermissionService{userRepo: userRepo}
}

func (s *PermissionService) HasPermission(userUUID uuid.UUID, permitType permission.Type) (bool, error) {
	u, err := s.userRepo.GetByUUID(userUUID)
	if err != nil {
		return false, err
	}
	if u == nil {
		return false, nil
	}

	return u.HasPermission(permitType), nil
}
