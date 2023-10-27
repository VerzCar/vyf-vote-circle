package repository

import (
	"database/sql"
	"fmt"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/api/model"
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
	Circles(userIdentityId string) ([]*model.Circle, error)
	UpdateCircle(circle *model.Circle) (*model.Circle, error)
	CreateNewCircle(circle *model.Circle) (*model.Circle, error)
	CountCirclesOfUser(userIdentityId string) (int64, error)

	CreateNewCircleVoter(voter *model.CircleVoter) (*model.CircleVoter, error)
	UpdateCircleVoter(voter *model.CircleVoter) (*model.CircleVoter, error)
	CircleVoterByCircleId(circleId int64, voterId string) (*model.CircleVoter, error)
	IsVoterInCircle(userIdentityId string, circle *model.Circle) (bool, error)

	RankingsByCircleId(circleId int64) ([]*model.Ranking, error)

	CreateNewVote(
		voterId int64,
		electedId int64,
		circleId int64,
	) (*model.Vote, error)
	ElectedVoterCountsByCircleId(circleId int64, electedId int64) (int64, error)
	VoterElectedByCircleId(
		circleId int64,
		voterId int64,
		electedId int64,
	) (*model.Vote, error)
	Votes(circleId int64) ([]*model.Vote, error)
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
