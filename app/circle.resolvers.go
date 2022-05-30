package app

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
)

func (r *mutationResolver) UpdateCircle(
	ctx context.Context,
	id int64,
	circleUpdateInput model.CircleUpdateInput,
) (*model.Circle, error) {
	circle, err := r.circleService.UpdateCircle(ctx, id, &circleUpdateInput)

	if err != nil {
		return nil, gqlerror.Errorf("cannot update circle")
	}

	return circle, nil
}

func (r *queryResolver) Circle(ctx context.Context, id int64) (*model.Circle, error) {
	circle, err := r.circleService.Circle(ctx, id)

	if err != nil {
		return nil, gqlerror.Errorf("cannot find circle")
	}

	return circle, nil
}
