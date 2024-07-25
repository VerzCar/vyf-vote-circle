package model

import (
	"github.com/VerzCar/vyf-vote-circle/app/database"
	_ "github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"time"
)

type RankingLastViewed struct {
	CreatedAt  time.Time `json:"createdAt" gorm:"autoCreateTime;"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"autoUpdateTime;"`
	Circle     *Circle   `json:"circle" gorm:"constraint:OnDelete:RESTRICT;"`
	IdentityID string    `json:"identityId" gorm:"type:varchar(50);not null"`
	ID         int64     `json:"id" gorm:"primary_key;"`
	CircleID   int64     `json:"circleId" gorm:"not null;"`
}

type RankingLastViewedResponse struct {
	CreatedAt time.Time                `json:"createdAt"`
	UpdatedAt time.Time                `json:"updatedAt"`
	Circle    *CirclePaginatedResponse `json:"circle"`
	ID        int64                    `json:"id"`
}

// db hooks with checks ++++++++++++++++++++++++++++

func (*RankingLastViewed) TableName() string {
	return "rankings_last_viewed"
}

func (ranking *RankingLastViewed) BeforeCreate(tx *gorm.DB) error {
	maxCountOfEntries := 15
	var exists bool

	err := tx.Model(&RankingLastViewed{}).
		Raw(
			"SELECT EXISTS(SELECT 1 FROM rankings_last_viewed WHERE identity_id=? AND circle_id=?)",
			ranking.IdentityID,
			ranking.CircleID,
		).
		Scan(&exists).
		Error

	if err != nil && !database.RecordNotFound(err) {
		return err
	}

	if exists {
		return DbErrEntryAlreadyExist
	}

	var rankings []*RankingLastViewed

	err = tx.Model(ranking).
		Where(&RankingLastViewed{IdentityID: ranking.IdentityID}).
		Order("created_at").
		Find(&rankings).Error

	if err != nil && !database.RecordNotFound(err) {
		return err
	}

	if len(rankings) >= maxCountOfEntries {
		err = tx.Model(ranking).
			Delete(&rankings[0]).Error

		if err != nil {
			return err
		}
	}

	return nil
}
