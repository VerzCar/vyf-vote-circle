package model

import (
	"time"
)

type CircleVoter struct {
	ID int64 `json:"id" gorm:"primary_key;"`
	// This must be a user identity id that should vote for the circle.
	Voter      string     `json:"voter" gorm:"type:varchar(50);not null"`
	Commitment Commitment `json:"commitment" gorm:"type:commitment;not null;default:OPEN"`
	// This must be a user identity id.
	VotedFor    *string   `json:"votedFor" gorm:"type:varchar(50)"`
	CircleID    int64     `json:"circleId" gorm:"not null;"`
	Circle      *Circle   `json:"circle" gorm:"constraint:OnDelete:RESTRICT"`
	CircleRefer *int64    `json:"circleRefer"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime;"`
}

type CircleVoterResponse struct {
	ID         int64      `json:"id"`
	Voter      string     `json:"voter"`
	Commitment Commitment `json:"commitment"`
	VotedFor   *string    `json:"votedFor"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

type CircleVotersResponse struct {
	Voters    []*CircleVoterResponse `json:"voters"`
	UserVoter *CircleVoterResponse   `json:"userVoter"`
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
