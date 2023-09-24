package app

/*
// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/graph/generated"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// RankingList is the resolver for the rankingList field.
func (r *queryResolver) RankingList(ctx context.Context, circleID int64) ([]*model.Ranking, error) {
	rankings, err := r.rankingService.Rankings(ctx, circleID)

	if err != nil {
		return nil, gqlerror.Errorf("cannot find rankings")
	}

	return rankings, nil
}

// RankingList is the resolver for the rankingList field.
func (r *subscriptionResolver) RankingList(ctx context.Context, circleID int64) (<-chan []*model.Ranking, error) {
	rankings, err := r.rankingSubscriptionService.Rankings(ctx, circleID)

	if err != nil {
		return nil, gqlerror.Errorf("cannot subscribe to ranking list for circle %d", circleID)
	}

	return rankings, nil
}

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
*/
