package services

import (
	"context"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/graph/client"
	appContext "gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/router/ctx"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/router/header"
)

type PaymentService interface {
	Payments(ctx context.Context) ([]*Payments, error)
}

type paymentService struct {
	host      string
	log       logger.Logger
	gqlClient *client.Client
}

func NewPaymentService(
	host string,
	log logger.Logger,
	gqlClient *client.Client,
) PaymentService {
	return &paymentService{
		host:      host,
		log:       log,
		gqlClient: gqlClient,
	}
}

type Payments struct {
	Id string
}

const (
	payments = `query payments {
	  payments {
		id
	  }
	}`
)

// Payments gets all the payements for the user
// Returns the payments for the user or an empty slice if no
// payments exists or an error if any occured.
func (p *paymentService) Payments(ctx context.Context) ([]*Payments, error) {
	var resp struct {
		Payments []*Payments
	}

	c, err := appContext.ContextToGinContext(ctx)

	if err != nil {
		p.log.Errorf("convert context failed: %s", err)
		return nil, err
	}

	accessToken, err := header.Authorization(c, "Bearer")

	if err != nil {
		p.log.Errorf("get accessToken from header failed: %s", err)
		return nil, err
	}

	err = p.gqlClient.Post(
		payments,
		&resp,
		client.Path(p.host),
		client.AddHeader(
			"Authorization",
			header.BearerToken(accessToken),
		),
	)

	if err != nil {
		p.log.Errorf("cannot get payments against payment svc: %s", err)
		return nil, err
	}

	return resp.Payments, nil
}
