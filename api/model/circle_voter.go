package model

import (
	"time"
)

type CircleVoter struct {
	CreatedAt   time.Time  `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"autoUpdateTime;"`
	VotedFor    *string    `json:"votedFor" gorm:"type:varchar(50)"`
	Circle      *Circle    `json:"circle" gorm:"constraint:OnDelete:RESTRICT"`
	CircleRefer *int64     `json:"circleRefer"`
	Voter       string     `json:"voter" gorm:"type:varchar(50);not null"`
	Commitment  Commitment `json:"commitment" gorm:"type:commitment;not null;default:OPEN"`
	ID          int64      `json:"id" gorm:"primary_key;"`
	CircleID    int64      `json:"circleId" gorm:"not null;"`
}

type CircleVoterResponse struct {
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	VotedFor   *string    `json:"votedFor"`
	Voter      string     `json:"voter"`
	Commitment Commitment `json:"commitment"`
	ID         int64      `json:"id"`
}

type CircleVotersResponse struct {
	UserVoter *CircleVoterResponse   `json:"userVoter"`
	Voters    []*CircleVoterResponse `json:"voters"`
}

type CircleVoterRequest struct {
	Voter string `json:"voter" validate:"gt=0,lte=50"`
}

type CircleVotersFilterBy struct {
	HasBeenVoted *bool `form:"hasBeenVoted,omitempty"`
}

type CircleVotersRequest struct {
	CircleVotersFilterBy
}

type CircleVoterChangedEvent struct {
	Voter     *CircleVoterResponse `json:"voter"`
	Operation EventOperation       `json:"operation"`
}
