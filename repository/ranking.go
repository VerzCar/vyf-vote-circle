package repository

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/database"
	"gorm.io/gorm"
)

// CreateNewRanking based on given model.Ranking model
func (s *storage) CreateNewRanking(ranking *model.Ranking) (*model.Ranking, error) {
	if err := s.db.Create(ranking).Error; err != nil {
		s.log.Errorf("error creating ranking: %s", err)
		return nil, err
	}

	return ranking, nil
}

func (s *storage) UpdateRanking(ranking *model.Ranking) (*model.Ranking, error) {
	if err := s.db.Model(&model.Ranking{ID: ranking.ID}).Updates(ranking).Error; err != nil {
		s.log.Errorf("error updating ranking: %s", err)
		return nil, err
	}

	return ranking, nil
}

// deletes ranking based on given ranking id
func (s *storage) DeleteRanking(rankingId int64) error {
	if err := s.db.Model(&model.Ranking{}).Delete(&model.Ranking{}, rankingId).Error; err != nil {
		s.log.Errorf("error deleting ranking: %s", err)
		return err
	}

	return nil
}

// RankingsByCircleId gets all rankings by the given circle id
func (s *storage) RankingsByCircleId(circleId int64) ([]*model.Ranking, error) {
	var rankings []*model.Ranking
	err := s.db.Where(&model.Ranking{CircleID: circleId}).
		Find(&rankings).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading rankings by circle id %d: %s", circleId, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("rankings with circle id %d not found: %s", circleId, err)
		return nil, err
	}

	return rankings, nil
}

func (s *storage) RankingByCircleId(circleId int64, identityId string) (*model.Ranking, error) {
	ranking := &model.Ranking{}
	err := s.db.Where(&model.Ranking{IdentityID: identityId, CircleID: circleId}).
		First(ranking).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading ranking by circle id %d for user %s: %s", circleId, identityId, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("ranking for user %s in circle %d not found: %s", identityId, circleId, err)
		return nil, err
	}

	return ranking, nil
}

func (s *storage) txUpsertRanking(
	tx *gorm.DB,
	circleId int64,
	voteCount int64,
	candidate *model.CircleCandidate,
) (*model.Ranking, error) {
	ranking := &model.Ranking{}

	err := tx.Where(&model.Ranking{IdentityID: candidate.Candidate, CircleID: circleId}).
		First(ranking).
		Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf(
			"error reading ranking by circle id %d for user %s: %s",
			circleId,
			candidate.Candidate,
			err,
		)
		return nil, err
	case database.RecordNotFound(err):
		newRanking := &model.Ranking{
			IdentityID: candidate.Candidate,
			Number:     0,
			Votes:      voteCount,
			CircleID:   circleId,
		}

		err = tx.Create(newRanking).Error

		if err != nil {
			s.log.Errorf("error creating ranking: %s", err)
			return nil, err
		}

		ranking = newRanking
		break
	default:
		ranking.Votes = voteCount

		err = tx.Model(ranking).
			Update("votes", voteCount).
			Error

		if err != nil {
			s.log.Errorf("error updating ranking: %s", err)
			return nil, err
		}
	}

	return ranking, err
}
