package auth

import (
	"errors"
	"net/http"

	"github.com/sak0/seeder/pkg/common/http/modifier"
	"github.com/sak0/seeder/pkg/common/secret"
)

// Authorizer is a kind of Modifier used to authorize the requests
type Authorizer modifier.Modifier

// SecretAuthorizer authorizes the requests with the specified secret
type SecretAuthorizer struct {
	secret string
}

// NewSecretAuthorizer returns an instance of SecretAuthorizer
func NewSecretAuthorizer(secret string) *SecretAuthorizer {
	return &SecretAuthorizer{
		secret: secret,
	}
}

// Modify the request by adding secret authentication information
func (s *SecretAuthorizer) Modify(req *http.Request) error {
	if req == nil {
		return errors.New("the request is null")
	}
	err := secret.AddToRequest(req, s.secret)
	return err
}
