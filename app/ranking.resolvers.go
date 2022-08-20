package app

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
)

// RankingList is the resolver for the rankingList field.
func (r *queryResolver) RankingList(ctx context.Context, circleID int64) ([]*model.Ranking, error) {
	panic(fmt.Errorf("not implemented: RankingList - rankingList"))
}
