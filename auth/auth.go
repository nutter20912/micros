package auth

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-micro/plugins/v4/auth/jwt"
	"github.com/go-micro/plugins/v4/auth/jwt/token"
	baseAuth "go-micro.dev/v4/auth"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/server"
)

var (
	publicKey  string
	privateKey string
	ttl        = time.Hour

	rules = []*baseAuth.Rule{
		// Enforce auth on all endpoints.
		{Scope: "*", Resource: &baseAuth.Resource{
			Type:     "*",
			Name:     "*",
			Endpoint: "*",
		}},

		// Enforce auth on one specific endpoint.
		{Scope: "*", Resource: &baseAuth.Resource{
			Type:     "service",
			Name:     "my.service.name",
			Endpoint: "UserService.Register",
		}},
	}
)

func init() {
	privateKeyBytes, err := os.ReadFile("private_key.pem")
	if err != nil {
		log.Println(err)
	} else {
		privateKey = string(base64.StdEncoding.EncodeToString(privateKeyBytes))
	}

	publicKeyBytes, err := os.ReadFile("public_key.pem")
	if err != nil {
		log.Println(err)
	} else {
		publicKey = string(base64.StdEncoding.EncodeToString(publicKeyBytes))
	}
}

func NewMicroAuth() baseAuth.Auth {
	return jwt.NewAuth(
		baseAuth.PublicKey(publicKey),
	)
}

func VerifyToken(a baseAuth.Auth, md metadata.Metadata) (*baseAuth.Account, error) {
	authHeader, ok := md["Authorization"]
	if !ok || !strings.HasPrefix(authHeader, baseAuth.BearerScheme) {
		return nil, errors.New("no auth token provided")
	}

	token := strings.TrimPrefix(authHeader, baseAuth.BearerScheme)

	acc, err := a.Inspect(token)
	if err != nil {
		return nil, errors.New("auth token invalid")
	}

	return acc, nil
}

func VerifyAbility(req server.Request, acc *baseAuth.Account) error {
	currentResource := baseAuth.Resource{
		Type:     "service",
		Name:     req.Service(),
		Endpoint: req.Endpoint(),
	}

	if err := baseAuth.Verify(rules, acc, &currentResource); err != nil {
		return errors.New("no access")
	}

	return nil
}

func Generate(id string) (*token.Token, error) {
	j := token.New(token.WithPrivateKey(privateKey))

	t, err := j.Generate(
		&baseAuth.Account{ID: id},
		token.WithExpiry(ttl),
	)
	if err != nil {
		return nil, fmt.Errorf("Generate returned %v error, expected nil", err)
	}

	return t, nil
}
