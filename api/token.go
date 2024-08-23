package api

import (
	"context"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/app/config"
	routerContext "github.com/VerzCar/vyf-vote-circle/app/router/ctx"
	"github.com/ably/ably-go/ably"
)

type TokenService interface {
	GenerateAblyToken(
		ctx context.Context,
	) (*ably.TokenRequest, error)
}

type tokenService struct {
	pubSubService *ably.Realtime
	config        *config.Config
	log           logger.Logger
}

func NewTokenService(
	pubSubService *ably.Realtime,
	config *config.Config,
	log logger.Logger,
) TokenService {
	return &tokenService{
		pubSubService: pubSubService,
		config:        config,
		log:           log,
	}
}

func (t *tokenService) GenerateAblyToken(
	ctx context.Context,
) (*ably.TokenRequest, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		t.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	params := &ably.TokenParams{
		ClientID: authClaims.PrivateClaims.ClientId,
	}

	tokenRequest, err := t.pubSubService.Auth.CreateTokenRequest(params)

	if err != nil {
		t.log.Errorf(
			"could not create token request: cause: %s",
			err,
		)
		return nil, err
	}

	return tokenRequest, nil
}
