package model

import (
	"database/sql/driver"
	"time"
)

type CircleCandidate struct {
	CreatedAt   time.Time  `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"autoUpdateTime;"`
	Circle      *Circle    `json:"circle" gorm:"constraint:OnDelete:RESTRICT"`
	CircleRefer *int64     `json:"circleRefer"`
	Candidate   string     `json:"voter" gorm:"type:varchar(50);not null"`
	Commitment  Commitment `json:"commitment" gorm:"type:commitment;not null;default:OPEN"`
	ID          int64      `json:"id" gorm:"primary_key;"`
	CircleID    int64      `json:"circleId" gorm:"not null;"`
}

type CircleCandidateResponse struct {
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	Candidate  string     `json:"candidate"`
	Commitment Commitment `json:"commitment"`
	ID         int64      `json:"id"`
}

type CircleCandidatesResponse struct {
	UserCandidate *CircleCandidateResponse   `json:"userCandidate"`
	Candidates    []*CircleCandidateResponse `json:"candidates"`
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

type CircleCandidateChangedEvent struct {
	Candidate *CircleCandidateResponse `json:"candidate"`
	Operation EventOperation           `json:"operation"`
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
