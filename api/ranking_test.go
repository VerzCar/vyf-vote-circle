package api_test

import (
	"context"
	"fmt"
	"github.com/VerzCar/vyf-vote-circle/api"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"reflect"
	"testing"
	"time"
)

func TestRankingService_Rankings(t *testing.T) {
	mockRepo := &mockRankingRepository{}
	mockCache := &mockRankingCache{}
	mockService := api.NewRankingService(mockRepo, mockCache, config, log)

	ctxMock := putUserIntoContext(mockUser.Elon)

	circleId := int64(1)
	circle02Id := int64(2)

	tests := []struct {
		name     string
		ctx      context.Context
		circleId int64
		want     error
	}{
		{
			name:     "should query rankings successfully",
			ctx:      ctxMock,
			circleId: circleId,
			want:     nil,
		},
		{
			name:     "should build rankings for circle successfully",
			ctx:      ctxMock,
			circleId: circle02Id,
			want:     nil,
		},
		{
			name:     "should fail because user is not authenticated",
			ctx:      emptyUserContext(),
			circleId: circleId,
			want:     fmt.Errorf("could not retrieve auth claims"),
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				_, err := mockService.Rankings(test.ctx, test.circleId)

				if !reflect.DeepEqual(err, test.want) {
					t.Errorf("test: %v failed. \ngot: %v \nwanted: %v", test.name, err, test.want)
				}
			},
		)
	}
}

type mockRankingRepository struct{}

type mockRankingCache struct{}

func (m mockRankingRepository) RankingsByCircleId(circleId int64) ([]*model.Ranking, error) {
	rankings := []*model.Ranking{
		{
			ID:         1,
			IdentityID: "1",
			Number:     1,
			Votes:      0,
			Placement:  "",
			CircleID:   circleId,
			Circle:     nil,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         2,
			IdentityID: "1",
			Number:     2,
			Votes:      0,
			Placement:  "",
			CircleID:   circleId,
			Circle:     nil,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}
	return rankings, nil
}

func (m mockRankingRepository) Votes(circleId int64) ([]*model.Vote, error) {
	votes := []*model.Vote{
		{
			ID:           1,
			VoterRefer:   1,
			Voter:        model.CircleVoter{},
			ElectedRefer: 2,
			Elected:      model.CircleVoter{},
			CircleID:     circleId,
			Circle:       nil,
			CircleRefer:  &circleId,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			ID:           2,
			VoterRefer:   2,
			Voter:        model.CircleVoter{},
			ElectedRefer: 3,
			Elected:      model.CircleVoter{},
			CircleID:     circleId,
			Circle:       nil,
			CircleRefer:  &circleId,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}
	return votes, nil
}

func (m mockRankingCache) RankingList(ctx context.Context, circleId int64) ([]*model.Ranking, error) {
	rankings := []*model.Ranking{
		{
			ID:         1,
			IdentityID: "1",
			Number:     1,
			Votes:      0,
			Placement:  "",
			CircleID:   circleId,
			Circle:     nil,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         2,
			IdentityID: "1",
			Number:     2,
			Votes:      0,
			Placement:  "",
			CircleID:   circleId,
			Circle:     nil,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}
	return rankings, nil
}

func (m mockRankingCache) ExistsRankingListForCircle(ctx context.Context, circleId int64) (bool, error) {
	switch circleId {
	case 2:
		return false, nil
	default:
		return true, nil
	}
}

func (m mockRankingCache) BuildRankingList(ctx context.Context, circleId int64, votes []*model.Vote) error {
	return nil
}
