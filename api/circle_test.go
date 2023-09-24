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

func TestCircleService_Circle(t *testing.T) {
	mockRepo := &mockCircleRepository{}
	mockService := api.NewCircleService(mockRepo, config, log)

	ctxMock := putUserIntoContext(mockUser.Elon)
	circleId := int64(1)
	circle02Id := int64(2)
	circle03Id := int64(3)

	tests := []struct {
		name     string
		ctx      context.Context
		circleId int64
		want     error
	}{
		{
			name:     "should query circle successfully",
			ctx:      ctxMock,
			circleId: circle02Id,
			want:     nil,
		},
		{
			name:     "should add user to global circle and query circle successfully",
			ctx:      ctxMock,
			circleId: circleId,
			want:     nil,
		},
		{
			name:     "should fail because use is not eligible to be in circle",
			ctx:      ctxMock,
			circleId: circle03Id,
			want:     fmt.Errorf("user is not eligible to be in circle"),
		},
		{
			name:     "should fail because circle does not exist",
			ctx:      ctxMock,
			circleId: -1,
			want:     fmt.Errorf("entry does not exist"),
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
				_, err := mockService.Circle(test.ctx, test.circleId)

				if !reflect.DeepEqual(err, test.want) {
					t.Errorf("test: %v failed. \ngot: %v \nwanted: %v", test.name, err, test.want)
				}
			},
		)
	}
}

func TestCircleService_CreateCircle(t *testing.T) {
	mockRepo := &mockCircleRepository{}
	mockService := api.NewCircleService(mockRepo, config, log)

	ctxMock := putUserIntoContext(mockUser.Elon)

	isPrivate := false
	description := "Beste election 2045"
	imageSrc := "https://source.com/img"
	validUntil := time.Now().Add(time.Hour * 8)
	circleCreateMockInput := model.CircleCreateRequest{
		Name:        "test1",
		Description: &description,
		ImageSrc:    &imageSrc,
		Voters: []*model.CircleVoterRequest{
			{
				Voter: mockUser.Elon.Subject,
			},
			{
				Voter: "voter02",
			},
		},
		Private:    &isPrivate,
		ValidUntil: &validUntil,
	}

	circleCreateMock02Input := circleCreateMockInput
	invalidValidUntilTime := time.Date(1990, time.Month(2), 21, 1, 10, 30, 0, time.UTC)
	circleCreateMock02Input.ValidUntil = &invalidValidUntilTime

	circleCreateMock03Input := circleCreateMockInput
	circleCreateMock03Input.Voters = make([]*model.CircleVoterRequest, 0)

	tests := []struct {
		name  string
		ctx   context.Context
		input *model.CircleCreateRequest
		want  error
	}{
		{
			name:  "should create circle successfully",
			ctx:   ctxMock,
			input: &circleCreateMockInput,
			want:  nil,
		},
		{
			name:  "should fail because given valid until time is in the past",
			ctx:   ctxMock,
			input: &circleCreateMock02Input,
			want:  fmt.Errorf("valid until time must be in the future from now"),
		},
		{
			name:  "should fail because voters are not given",
			ctx:   ctxMock,
			input: &circleCreateMock03Input,
			want:  fmt.Errorf("voters for circle are not given"),
		},
		{
			name:  "should fail because user is not authenticated",
			ctx:   emptyUserContext(),
			input: &circleCreateMockInput,
			want:  fmt.Errorf("could not retrieve auth claims"),
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				_, err := mockService.CreateCircle(test.ctx, test.input)

				if !reflect.DeepEqual(err, test.want) {
					t.Errorf("test: %v failed. \ngot: %v \nwanted: %v", test.name, err, test.want)
				}
			},
		)
	}
}

func TestCircleService_UpdateCircle(t *testing.T) {
	mockRepo := &mockCircleRepository{}
	mockService := api.NewCircleService(mockRepo, config, log)

	ctxMock := putUserIntoContext(mockUser.Elon)

	circleId := int64(1)
	circle02Id := int64(2)
	circle04Id := int64(4)

	name := "test1"
	isPrivate := true
	description := "Best second election 2045"
	imageSrc := "https://source.com/img/2"
	validUntil := time.Now().Add(time.Hour * 10)
	deleteCircle := false
	circleUpdateMockInput := model.CircleUpdateInput{
		Name:        &name,
		Description: &description,
		ImageSrc:    &imageSrc,
		Voters: []*model.CircleVoterRequest{
			{
				Voter: mockUser.Elon.Subject,
			},
			{
				Voter: "voter02",
			},
		},
		Private:    &isPrivate,
		ValidUntil: &validUntil,
		Delete:     &deleteCircle,
	}

	circleUpdateMock02Input := circleUpdateMockInput
	deleteCircle02 := true
	circleUpdateMock02Input.Delete = &deleteCircle02

	tests := []struct {
		name     string
		ctx      context.Context
		circleId int64
		input    *model.CircleUpdateInput
		want     error
	}{
		{
			name:     "should update circle successfully",
			ctx:      ctxMock,
			circleId: circle02Id,
			input:    &circleUpdateMockInput,
			want:     nil,
		},
		{
			name:     "should update (delete) circle successfully",
			ctx:      ctxMock,
			circleId: circle02Id,
			input:    &circleUpdateMock02Input,
			want:     nil,
		},
		{
			name:     "should fail because circle is not active anymore",
			ctx:      ctxMock,
			circleId: circle04Id,
			input:    &circleUpdateMockInput,
			want:     fmt.Errorf("circle is not active"),
		},
		{
			name:     "should fail because user is not eligible to update circle",
			ctx:      ctxMock,
			circleId: circleId,
			input:    &circleUpdateMockInput,
			want:     fmt.Errorf("user is not eligible to update circle"),
		},
		{
			name:     "should fail because user is not authenticated",
			ctx:      emptyUserContext(),
			circleId: circleId,
			input:    &circleUpdateMockInput,
			want:     fmt.Errorf("could not retrieve auth claims"),
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				_, err := mockService.UpdateCircle(test.ctx, test.circleId, test.input)

				if !reflect.DeepEqual(err, test.want) {
					t.Errorf("test: %v failed. \ngot: %v \nwanted: %v", test.name, err, test.want)
				}
			},
		)
	}
}

type mockCircleRepository struct{}

func (m mockCircleRepository) UpdateCircle(circle *model.Circle) (*model.Circle, error) {
	return circle, nil
}

func (m mockCircleRepository) CreateNewCircle(circle *model.Circle) (*model.Circle, error) {
	return circle, nil
}

func (m mockCircleRepository) CreateNewCircleVoter(voter *model.CircleVoter) (*model.CircleVoter, error) {
	newVoter := &model.CircleVoter{
		ID:          1,
		Voter:       "test1",
		Commitment:  "",
		CircleID:    0,
		Circle:      nil,
		CircleRefer: nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return newVoter, nil
}

func (m mockCircleRepository) CircleById(id int64) (*model.Circle, error) {
	circle := &model.Circle{
		ID:          id,
		Name:        "Circle 1",
		Description: "Circle 1 description",
		ImageSrc:    "",
		Votes:       []*model.Vote{},
		Voters:      []*model.CircleVoter{},
		Private:     false,
		Active:      true,
		CreatedFrom: "anonymous",
		ValidUntil:  nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	switch id {
	case 2:
		circle.CreatedFrom = mockUser.Elon.Subject
		break
	case 4:
		circle.CreatedFrom = mockUser.Elon.Subject
		circle.Active = false
		break
	case -1:
		return nil, fmt.Errorf("entry does not exist")
	}

	return circle, nil
}
