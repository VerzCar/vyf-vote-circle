package app

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
)

// RankingList is the resolver for the rankingList field.
func (r *queryResolver) RankingList(ctx context.Context, circleID int64) ([]*model.Ranking, error) {
	rankings, err := r.rankingService.Rankings(ctx, circleID)

	if err != nil {
		return nil, gqlerror.Errorf("cannot find rankings")
	}

	return rankings, nil
}