package model

import (
	"database/sql/driver"
	"encoding/json"
	_ "github.com/go-playground/validator/v10"
	"time"
)

type Ranking struct {
	CreatedAt  time.Time `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"autoUpdateTime;"`
	Circle     *Circle   `json:"circle" gorm:"constraint:OnDelete:RESTRICT;"`
	IdentityID string    `json:"identityId" gorm:"type:varchar(50);not null"`
	Placement  Placement `json:"placement" gorm:"type:placement;not null;default:NEUTRAL"`
	ID         int64     `json:"id" gorm:"primary_key;index;"`
	Number     int64     `json:"number" gorm:"primary_key;index;"`
	Votes      int64     `json:"votes" gorm:"not null;default:0"`
	CircleID   int64     `json:"circleId" gorm:"not null;"`
}

type RankingResponse struct {
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	IdentityID   string    `json:"identityId"`
	Placement    Placement `json:"placement"`
	ID           int64     `json:"id"`
	CandidateID  int64     `json:"candidateId"`
	Number       int64     `json:"number"`
	Votes        int64     `json:"votes"`
	IndexedOrder int64     `json:"indexedOrder"`
	CircleID     int64     `json:"circleId"`
}

type RankingsUriRequest struct {
	CircleID int64 `uri:"circleId" validate:"gt=0"`
}

type RankingScore struct {
	UserIdentityId string `redis:"userIdentityId"`
	VoteCount      int64  `redis:"voteCount"`
}

type RankingUserCandidate struct {
	CreatedAt   time.Time `redis:"time"`
	UpdatedAt   time.Time `redis:"time"`
	CandidateID int64     `redis:"candidateId"`
	RankingID   int64     `redis:"rankingId"`
}

type RankingCacheItem struct {
	Ranking   *Ranking
	Candidate *CircleCandidate
	VoteCount int64
}

type RankingChangedEvent struct {
	Ranking   *RankingResponse `json:"ranking"`
	Operation EventOperation   `json:"operation"`
}

func (s RankingScore) MarshalBinary() ([]byte, error) {
	return json.Marshal(s)
}

func (s RankingScore) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &s)
}

type Placement string

const (
	PlacementNeutral    Placement = "NEUTRAL"
	PlacementAscending  Placement = "ASCENDING"
	PlacementDescending Placement = "DESCENDING"
)

var AllPlacement = []Placement{
	PlacementNeutral,
	PlacementAscending,
	PlacementDescending,
}

func (e *Placement) Scan(value interface{}) error {
	*e = Placement(value.(string))
	return nil
}

func (e Placement) Value() (driver.Value, error) {
	return string(e), nil
}

func (e Placement) IsValid() bool {
	switch e {
	case PlacementNeutral, PlacementAscending, PlacementDescending:
		return true
	}
	return false
}

func (e Placement) String() string {
	return string(e)
}
