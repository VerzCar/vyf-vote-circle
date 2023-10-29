package model

import (
	"database/sql/driver"
	"encoding/json"
	_ "github.com/go-playground/validator/v10"
	"time"
)

type Ranking struct {
	ID         int64     `json:"id" gorm:"primary_key;index;"`
	IdentityID string    `json:"identityId" gorm:"type:varchar(50);not null"`
	Number     int64     `json:"number" gorm:"primary_key;index;"`
	Votes      int64     `json:"votes" gorm:"not null;default:0"`
	Placement  Placement `json:"placement" gorm:"type:placement;not null;default:NEUTRAL"`
	CircleID   int64     `json:"circleId" gorm:"not null;"`
	Circle     *Circle   `json:"circle" gorm:"constraint:OnDelete:RESTRICT;"`
	CreatedAt  time.Time `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"autoUpdateTime;"`
}

type RankingResponse struct {
	ID         int64     `json:"id"`
	IdentityID string    `json:"identityId"`
	Number     int64     `json:"number"`
	Votes      int64     `json:"votes"`
	Placement  Placement `json:"placement"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type RankingsUriRequest struct {
	CircleID int64 `uri:"circleId" validate:"gt=0"`
}

type RankingScore struct {
	VoteCount      int64
	UserIdentityId string
}

type RankingListObservable chan []*Ranking
type RankingListObservableMap map[string]RankingListObservable
type RankingObservers map[int64]RankingListObservableMap

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
