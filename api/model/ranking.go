package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	_ "github.com/go-playground/validator/v10"
	"io"
	"strconv"
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

type RankingsUriRequest struct {
	CircleID int64 `uri:"circleId" validate:"gt=0"`
}

type RankingScore struct {
	VoteCount      int64
	UserIdentityId string
}

type RankingObservers map[int64]RankingListObservableMap
type RankingListObservableMap map[string]RankingListObservable
type RankingListObservable chan []*Ranking

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

func (e *Placement) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Placement(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Placement", str)
	}
	return nil
}

func (e Placement) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
