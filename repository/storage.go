package repository

import (
	"context"
	"database/sql"
	"fmt"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/cache"
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/VerzCar/vyf-vote-circle/app/database"
	"github.com/VerzCar/vyf-vote-circle/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"path/filepath"
)

type Storage interface {
	RunMigrationsUp(db *sql.DB) error
	RunMigrationsDown(db *sql.DB) error

	CircleById(id int64) (*model.Circle, error)
	CirclesByIds(circleIds []int64) ([]*model.CirclePaginated, error)
	Circles(userIdentityId string) ([]*model.Circle, error)
	CirclesFiltered(name string) ([]*model.CirclePaginated, error)
	CirclesOfInterest(userIdentityId string) ([]*model.CirclePaginated, error)
	UpdateCircle(circle *model.Circle) (*model.Circle, error)
	CreateNewCircle(circle *model.Circle) (*model.Circle, error)
	CountCirclesOfUser(userIdentityId string) (int64, error)

	CreateNewCircleVoter(voter *model.CircleVoter) (*model.CircleVoter, error)
	UpdateCircleVoter(voter *model.CircleVoter) (*model.CircleVoter, error)
	DeleteCircleVoter(voterId int64) error
	CircleVoterByCircleId(circleId int64, userIdentityId string) (*model.CircleVoter, error)
	CircleVoterCountByCircleId(
		circleId int64,
	) (int64, error)
	IsVoterInCircle(userIdentityId string, circleId int64) (bool, error)
	CircleVotersFiltered(
		circleId int64,
		filterBy *model.CircleVotersFilterBy,
	) ([]*model.CircleVoter, error)
	CircleVotersVotedFor(
		circleId int64,
		userIdentityId string,
	) ([]*model.CircleVoter, error)

	CreateNewCircleCandidate(candidate *model.CircleCandidate) (*model.CircleCandidate, error)
	UpdateCircleCandidate(candidate *model.CircleCandidate) (*model.CircleCandidate, error)
	DeleteCircleCandidate(candidateId int64) error
	CircleCandidateByCircleId(
		circleId int64,
		userIdentityId string,
	) (*model.CircleCandidate, error)
	IsCandidateInCircle(
		userIdentityId string,
		circleId int64,
	) (bool, error)
	CircleCandidateCountByCircleId(
		circleId int64,
	) (int64, error)
	CircleCandidatesFiltered(
		circleId int64,
		filterBy *model.CircleCandidatesFilterBy,
	) ([]*model.CircleCandidate, error)
	CircleCandidatesOpenCommitments(
		userIdentityId string,
	) ([]*model.CircleCandidate, error)

	CreateNewRanking(ranking *model.Ranking) (*model.Ranking, error)
	UpdateRanking(ranking *model.Ranking) (*model.Ranking, error)
	DeleteRanking(rankingId int64) error
	RankingsByCircleId(circleId int64) ([]*model.Ranking, error)
	RankingByCircleId(circleId int64, identityId string) (*model.Ranking, error)

	CreateNewVote(
		ctx context.Context,
		circleId int64,
		voter *model.CircleVoter,
		candidate *model.CircleCandidate,
		upsertRankingCache cache.UpsertRankingCacheCallback,
	) (*model.RankingResponse, int64, error)
	DeleteVote(
		ctx context.Context,
		circleId int64,
		vote *model.Vote,
		voter *model.CircleVoter,
		upsertRankingCache cache.UpsertRankingCacheCallback,
		removeRankingCache cache.RemoveRankingCacheCallback,
	) (*model.RankingResponse, int64, error)
	VoteByCircleId(
		circleId int64,
		voterId int64,
	) (*model.Vote, error)
	CountsVotesOfCandidateByCircleId(circleId int64, candidateId int64) (int64, error)
	HasVoterVotedForCircle(
		circleId int64,
		voterId int64,
	) (bool, error)
	Votes(circleId int64) ([]*model.Vote, error)

	CreateNewUserOption(option *model.UserOption) (*model.UserOption, error)
	DeleteUserOption(optionId int64) error
	UserOptionByUserIdentityId(userIdentityId string) (*model.UserOption, error)
}

type storage struct {
	db     database.Client
	config *config.Config
	log    logger.Logger
}

func NewStorage(
	db database.Client,
	config *config.Config,
	log logger.Logger,
) Storage {
	return &storage{
		db:     db,
		config: config,
		log:    log,
	}
}

func (s *storage) RunMigrationsUp(db *sql.DB) error {
	m, err := createMigrateInstance(db)

	if err != nil {
		return err
	}

	err = m.Up()

	switch err {
	case migrate.ErrNoChange:
		return nil
	}

	return err
}

func (s *storage) RunMigrationsDown(db *sql.DB) error {
	m, err := createMigrateInstance(db)

	if err != nil {
		return err
	}

	err = m.Down()

	switch err {
	case migrate.ErrNoChange:
		return nil
	}

	return err
}

func createMigrateInstance(db *sql.DB) (*migrate.Migrate, error) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})

	if err != nil {
		return nil, fmt.Errorf("error creating migrations with db instance: %s", err)
	}

	repoMigrationPath := utils.FromBase("repository/migrations/")
	migrationsPath := filepath.Join("file://", repoMigrationPath)

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres", driver,
	)

	if err != nil {
		return nil, err
	}

	return m, nil
}
