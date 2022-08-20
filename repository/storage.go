package repository

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/config"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/utils"
	"gorm.io/gorm"
	"path/filepath"
)

type Storage interface {
	RunMigrationsUp(db *sql.DB) error
	RunMigrationsDown(db *sql.DB) error
	CircleById(id int64) (*model.Circle, error)
	UpdateCircle(circle *model.Circle) (*model.Circle, error)
	CreateNewCircle(circle *model.Circle) (*model.Circle, error)

	CreateNewCircleVoter(voter *model.CircleVoter) (*model.CircleVoter, error)
}

type storage struct {
	db     *gorm.DB
	config *config.Config
	log    logger.Logger
}

func NewStorage(
	db *gorm.DB,
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
