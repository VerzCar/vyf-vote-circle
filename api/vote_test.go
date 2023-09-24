package api_test

import (
	"context"
	"fmt"
	"github.com/VerzCar/vyf-vote-circle/api"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"gorm.io/gorm"
	"reflect"
	"testing"
	"time"
)

func TestVoteService_Vote(t *testing.T) {
	mockRepo := &mockVoteRepository{}
	mockCache := &mockVoteCache{}
	mockVoteSubscriptionSvc := &mockVoteSubscriptionService{}
	mockService := api.NewVoteService(mockRepo, mockCache, mockVoteSubscriptionSvc, config, log)

	ctxMock := putUserIntoContext(mockUser.Elon)
	voteInput := &model.VoteCreateInput{
		Elected: "test1",
	}
	circleId := int64(1)

	inactiveCircleId := int64(2)

	tests := []struct {
		name     string
		ctx      context.Context
		circleId int64
		input    *model.VoteCreateInput
		want     error
	}{
		{
			name:     "should create a new vote successfully",
			ctx:      ctxMock,
			circleId: circleId,
			input:    voteInput,
			want:     nil,
		},
		{
			name:     "should fail because circle is inactive",
			ctx:      ctxMock,
			circleId: inactiveCircleId,
			input:    voteInput,
			want:     fmt.Errorf("circle inactive"),
		},
		{
			name:     "should fail because user is not authenticated",
			ctx:      emptyUserContext(),
			circleId: inactiveCircleId,
			input:    voteInput,
			want:     fmt.Errorf("could not retrieve auth claims"),
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				_, err := mockService.CreateVote(test.ctx, test.circleId, test.input)

				if !reflect.DeepEqual(err, test.want) {
					t.Errorf("test: %v failed. \ngot: %v \nwanted: %v", test.name, err, test.want)
				}
			},
		)
	}
}

type mockVoteRepository struct{}

type mockVoteCache struct{}

type mockVoteSubscriptionService struct{}

func (m mockVoteRepository) CircleById(id int64) (*model.Circle, error) {
	circle := &model.Circle{
		ID:          id,
		Name:        "Circle 1",
		Description: "Circle 1 description",
		ImageSrc:    "",
		Votes:       nil,
		Voters:      nil,
		Private:     false,
		Active:      true,
		CreatedFrom: "",
		ValidUntil:  nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	switch id {
	case 1:
		break
	case 2:
		circle.Active = false
		break
	}

	return circle, nil
}

func (m mockVoteRepository) CircleVoterByCircleId(circleId int64, voterId string) (*model.CircleVoter, error) {
	voter := &model.CircleVoter{
		ID:          1,
		Voter:       "test1",
		Commitment:  "",
		CircleID:    circleId,
		Circle:      nil,
		CircleRefer: &circleId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return voter, nil
}

func (m mockVoteRepository) CreateNewVote(
	voterId int64,
	electedId int64,
	circleId int64,
) (*model.Vote, error) {
	vote := &model.Vote{
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
	}
	return vote, nil
}

func (m mockVoteRepository) ElectedVoterCountsByCircleId(circleId int64, electedId int64) (int64, error) {
	return 3, nil
}

func (m mockVoteRepository) VoterElectedByCircleId(
	circleId int64,
	voterId int64,
	electedId int64,
) (*model.Vote, error) {
	return nil, gorm.ErrRecordNotFound
}

func (m mockVoteCache) UpdateRanking(
	ctx context.Context,
	circleId int64,
	identityId string,
	votes int64,
) error {
	return nil
}

func (m mockVoteSubscriptionService) RankingChangedEvent(circleId int64) {
	return
}
