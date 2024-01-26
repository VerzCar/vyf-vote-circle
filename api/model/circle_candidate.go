package model

import (
	"database/sql/driver"
	"time"
)

type CircleCandidate struct {
	ID int64 `json:"id" gorm:"primary_key;"`
	// This must be a user identity id that should be a candidate for the circle.
	Candidate   string     `json:"voter" gorm:"type:varchar(50);not null"`
	Commitment  Commitment `json:"commitment" gorm:"type:commitment;not null;default:OPEN"`
	CircleID    int64      `json:"circleId" gorm:"not null;"`
	Circle      *Circle    `json:"circle" gorm:"constraint:OnDelete:RESTRICT"`
	CircleRefer *int64     `json:"circleRefer"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"autoUpdateTime;"`
}

type CircleCandidateResponse struct {
	ID         int64      `json:"id"`
	Candidate  string     `json:"candidate"`
	Commitment Commitment `json:"commitment"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

type CircleCandidatesResponse struct {
	Candidates    []*CircleCandidateResponse `json:"candidates"`
	UserCandidate *CircleCandidateResponse   `json:"userCandidate"`
}

type CircleCandidateRequest struct {
	Candidate string `json:"candidate" validate:"gt=0,lte=50"`
}

type CircleCandidateCommitmentRequest struct {
	Commitment Commitment `json:"commitment" validate:"gt=0,lte=20"`
}

type CircleCandidatesFilterBy struct {
	Commitment   *Commitment `form:"commitment,omitempty" validate:"omitempty,gt=0,lte=12"`
	HasBeenVoted *bool       `form:"hasBeenVoted,omitempty"`
}

type CircleCandidatesRequest struct {
	CircleCandidatesFilterBy
}

type Commitment string

const (
	CommitmentOpen      Commitment = "OPEN"
	CommitmentCommitted Commitment = "COMMITTED"
	CommitmentRejected  Commitment = "REJECTED"
)

var AllCommitment = []Commitment{
	CommitmentOpen,
	CommitmentCommitted,
	CommitmentRejected,
}

func (e *Commitment) Scan(value interface{}) error {
	*e = Commitment(value.(string))
	return nil
}

func (e Commitment) Value() (driver.Value, error) {
	return string(e), nil
}

func (e Commitment) IsValid() bool {
	switch e {
	case CommitmentOpen, CommitmentCommitted, CommitmentRejected:
		return true
	}
	return false
}

func (e Commitment) String() string {
	return string(e)
}
