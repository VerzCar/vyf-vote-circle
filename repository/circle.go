package repository

import (
	"fmt"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CircleById gets the circle by id
func (s *storage) CircleById(id int64) (*model.Circle, error) {
	circle := &model.Circle{}
	err := s.db.Where(&model.Circle{ID: id, Active: true}).
		First(circle).
		Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading circle by id %d: %s", id, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("circle with id %d not found: %s", id, err)
		return nil, err
	}

	return circle, nil
}

// Circles gets all the active circles that have been created from the user
func (s *storage) Circles(userIdentityId string) ([]*model.Circle, error) {
	var circles []*model.Circle
	err := s.db.Where(&model.Circle{CreatedFrom: userIdentityId, Active: true}).
		Limit(int(s.config.Circle.MaxAmountPerUser)).
		Order("updated_at desc").
		Find(&circles).
		Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading circles for user id %s: %s", userIdentityId, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("circles with user id %s not found: %s", userIdentityId, err)
		return nil, err
	}

	return circles, nil
}

// CirclesFiltered gets all active circles that matches the filter
func (s *storage) CirclesFiltered(name string) ([]*model.CirclePaginated, error) {
	var circles []*model.CirclePaginated

	err := s.db.Model(&model.Circle{}).
		Select("circles.id, circles.name, circles.description, circles.image_src, circles.active, circles.created_at, circles.updated_at").
		Where("name LIKE ?", fmt.Sprintf("%%%s%%", name)).
		Where(&model.Circle{Active: true}).
		Limit(100).
		Order("updated_at desc").
		Find(&circles).
		Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading circles: %s", err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("circles not found: %s", err)
		return nil, err
	}

	return circles, nil
}

// CirclesOfInterest evaluates all the circles that the user is involved (is a voter)
// and filters out the ones that belongs to the user.
func (s *storage) CirclesOfInterest(userIdentityId string) ([]*model.CirclePaginated, error) {
	rows, err := s.db.Model(&model.Circle{}).Raw(
		`SELECT *
			FROM (SELECT Distinct ON (circles.id) circles.id,
												  circles.name,
												  circles.description,
												  circles.image_src,
												  circles.active,
												  circles.created_at,
												  circles.updated_at
            FROM circles
                 left join circle_voters voters on circles.id = voters.circle_id
                 left join circle_candidates candidates on circles.id = candidates.circle_id
			  WHERE circles.active = ?
                 AND circles.created_from <> ?
                 AND (voters.circle_id IS NULL
	               OR ((circles.private = ? AND voters.voter = ?)
	               OR (circles.private = ? AND voters.voter <> ?)))
			     OR (candidates.circle_id IS NULL
	               OR ((circles.private = ? AND candidates.candidate = ?)
	               OR (circles.private = ? AND candidates.candidate <> ?)))
			  LIMIT ?) circles_of_interest
            ORDER BY updated_at desc;`,
		true,
		userIdentityId,
		true,
		userIdentityId,
		false,
		userIdentityId,
		true,
		userIdentityId,
		false,
		userIdentityId,
		100,
	).Rows()

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading nearest circles: %s", err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("nearest circles not found: %s", err)
		return nil, err
	}

	circles := make([]*model.CirclePaginated, 0)

	defer rows.Close()
	for rows.Next() {
		circle := &model.CirclePaginated{}
		err := rows.Scan(
			&circle.ID,
			&circle.Name,
			&circle.Description,
			&circle.ImageSrc,
			&circle.Active,
			&circle.CreatedAt,
			&circle.UpdatedAt,
		)

		switch {
		case err != nil && !database.RecordNotFound(err):
			s.log.Errorf("error reading nearest circles: %s", err)
			return nil, err
		case database.RecordNotFound(err):
			s.log.Infof("nearest circles not found: %s", err)
			return nil, err
		}

		err = s.db.Model(&model.Circle{}).Raw(
			`	SELECT count(1) as voters_count
	             from circle_voters
	             WHERE circle_id = ?
	             group by circle_voters.circle_id;`,
			circle.ID,
		).Scan(&circle.VotersCount).Error

		switch {
		case err != nil && !database.RecordNotFound(err):
			s.log.Errorf("error reading count of voters for circles: %s", err)
			return nil, err
		case database.RecordNotFound(err):
			s.log.Infof("counts of voters for circle not found: %s", err)
			return nil, err
		}

		err = s.db.Model(&model.Circle{}).Raw(
			`	SELECT count(1) as candidates_count
	             from circle_candidates
	             WHERE circle_id = ?
	             group by circle_candidates.circle_id;`,
			circle.ID,
		).Scan(&circle.CandidatesCount).Error

		switch {
		case err != nil && !database.RecordNotFound(err):
			s.log.Errorf("error reading count of candidates for circles: %s", err)
			return nil, err
		case database.RecordNotFound(err):
			s.log.Infof("counts of candidates for circle not found: %s", err)
			return nil, err
		}

		circles = append(circles, circle)
	}

	return circles, nil
}

// UpdateCircle update circle based on given circle model
func (s *storage) UpdateCircle(circle *model.Circle) (*model.Circle, error) {
	if err := s.db.Save(circle).Error; err != nil {
		s.log.Errorf("error updating circle: %s", err)
		return nil, err
	}

	return circle, nil
}

// CreateNewCircle based on given circle model.
// The associations that come with it, will be created in the transaction accordingly.
func (s *storage) CreateNewCircle(circle *model.Circle) (*model.Circle, error) {
	err := s.db.Transaction(
		func(tx *gorm.DB) error {
			err := tx.Model(circle).Omit(clause.Associations).Create(circle).Error

			if err != nil {
				s.log.Error("error creating circle entry: %s", err)
				return err
			}

			circleVoters := circle.Voters

			if len(circleVoters) > 0 {
				for _, voter := range circleVoters {
					voter.CircleID = circle.ID
					voter.CircleRefer = &circle.ID
				}

				err = tx.Model(&model.CircleVoter{}).Create(circleVoters).Error

				if err != nil {
					s.log.Error("error creating circle voters entry: %s", err)
					return err
				}
			}

			circle.Voters = circleVoters

			circleCandidates := circle.Candidates

			if len(circleCandidates) > 0 {
				for _, candidate := range circleCandidates {
					candidate.CircleID = circle.ID
					candidate.CircleRefer = &circle.ID
				}

				err = tx.Model(&model.CircleCandidate{}).Create(circleCandidates).Error

				if err != nil {
					s.log.Error("error creating circle candidates entry: %s", err)
					return err
				}
			}

			circle.Candidates = circleCandidates
			return nil
		},
	)

	if err != nil {
		s.log.Error("error creating circle: %s", err)
		return nil, err
	}

	return circle, nil
}

// CountCirclesOfUser determines how many circles the user already obtains
func (s *storage) CountCirclesOfUser(
	userIdentityId string,
) (int64, error) {
	var count int64
	err := s.db.Model(&model.Circle{}).
		Where(&model.Circle{CreatedFrom: userIdentityId, Active: true}).
		Count(&count).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading circle count by user id %s: %s", userIdentityId, err)
		return 0, err
	case database.RecordNotFound(err):
		s.log.Infof("user with id %s in circles not found: %s", userIdentityId, err)
		return 0, err
	}

	return count, nil
}
