package idp

import (
	"context"

	"github.com/Nerzal/gocloak/v12"
	"github.com/google/uuid"

	"github.com/jbenzshawel/go-sandbox/identity/domain/user"
)

type IdentityProvider interface {
	CreateUser(ctx context.Context, user *user.User, password string) (uuid.UUID, error)
	UpdateUser(ctx context.Context, user *user.User) error
	DeleteUser(ctx context.Context, userID string) error
}

type KeyCloakProvider struct {
	client   *gocloak.GoCloak
	user     string
	password string
	realm    string
}

func NewKeyCloakProvider(basePath, user, password, realm string) *KeyCloakProvider {
	return &KeyCloakProvider{
		client:   gocloak.NewClient(basePath),
		user:     user,
		password: password,
		realm:    realm,
	}
}

func (idp *KeyCloakProvider) CreateUser(ctx context.Context, user *user.User, password string) (uuid.UUID, error) {
	jwt, err := idp.getToken(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	idpUser := idp.mapKeyCloakUser(user)

	idpUserID, err := idp.client.CreateUser(ctx, jwt.AccessToken, idp.realm, idpUser)
	if err != nil {
		return uuid.Nil, err
	}

	userUUID, err := uuid.Parse(idpUserID)
	if err != nil {
		return uuid.Nil, err
	}

	err = idp.client.SetPassword(ctx, jwt.AccessToken, idpUserID, idp.realm, password, false)
	if err != nil {
		return uuid.Nil, err
	}

	return userUUID, nil
}

func (idp *KeyCloakProvider) UpdateUser(ctx context.Context, user *user.User) error {
	jwt, err := idp.getToken(ctx)
	if err != nil {
		return err
	}

	idpUser := idp.mapKeyCloakUser(user)
	idpUser.ID = gocloak.StringP(user.UUID().String())

	return idp.client.UpdateUser(ctx, jwt.AccessToken, idp.realm, idpUser)
}

func (idp *KeyCloakProvider) mapKeyCloakUser(user *user.User) gocloak.User {
	return gocloak.User{
		FirstName:     gocloak.StringP(user.FirstName()),
		LastName:      gocloak.StringP(user.LastName()),
		Email:         gocloak.StringP(user.Email()),
		EmailVerified: gocloak.BoolP(user.EmailVerified()),
		Enabled:       gocloak.BoolP(user.Enabled()),
		Username:      gocloak.StringP(user.Email()),
	}
}

func (idp *KeyCloakProvider) DeleteUser(ctx context.Context, userID string) error {
	jwt, err := idp.getToken(ctx)
	if err != nil {
		return err
	}

	return idp.client.DeleteUser(ctx, jwt.AccessToken, idp.realm, userID)
}

func (idp *KeyCloakProvider) getToken(ctx context.Context) (*gocloak.JWT, error) {
	return idp.client.LoginAdmin(ctx, idp.user, idp.password, "master")
}
