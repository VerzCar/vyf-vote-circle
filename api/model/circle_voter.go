package model

import (
	"database/sql/driver"
	"fmt"
	"io"
	"strconv"
	"time"
)

type CircleVoter struct {
	ID int64 `json:"id" gorm:"primary_key;index;"`

	// Voter that is eligible to vote.
	// This must be a user identity id that should vote for this circle.
	Voter string `json:"voter" gorm:"type:varchar(50);not null"`

	Commitment Commitment `json:"commitment" gorm:"type:commitment;not null;default:OPEN"`

	CircleID    int64     `json:"circleId" gorm:"not null;"`
	Circle      *Circle   `json:"circle" gorm:"constraint:OnDelete:RESTRICT"`
	CircleRefer *int64    `json:"circleRefer"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime;"`
}

type CircleVoterInput struct {
	Voter string `json:"voter" validate:"omitempty,gt=0,lte=50"`
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

func (e *Commitment) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Commitment(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Commitment", str)
	}
	return nil
}

func (e Commitment) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
