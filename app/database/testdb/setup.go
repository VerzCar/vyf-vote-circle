package testdb

import (
	"fmt"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/config"
	"gorm.io/gorm"
)

// Setup will setup the test database.
// It creates the database and migrates the tables from the models.
func Setup(db *gorm.DB, log logger.Logger, conf *config.Config) {
	createDb(db, log, conf)
}

// createDb will create the test database.
func createDb(db *gorm.DB, log logger.Logger, conf *config.Config) {
	createDbStatement := fmt.Sprintf(
		`SELECT 'CREATE DATABASE %s' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s')`,
		conf.Db.Test.Name, conf.Db.Test.Name,
	)

	if err := db.Exec(createDbStatement).Error; err != nil {
		log.Fatalf("could not create database %s. error: %s", conf.Db.Test.Name, err)
	}
}
