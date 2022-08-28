package app

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
)

// CreateVote is the resolver for the createVote field.
func (r *mutationResolver) CreateVote(ctx context.Context, circleID int64, voteCreateInput model.VoteCreateInput) (bool, error) {
	result, err := r.voteService.Vote(ctx, circleID, &voteCreateInput)

	if err != nil {
		return false, gqlerror.Errorf("cannot vote")
	}

	return result, nil
}
