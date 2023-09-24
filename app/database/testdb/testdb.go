package testdb

import (
	"fmt"
	"github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect to the database with configured dsn.
// If successful, the gorm.DB connection will be returned, otherwise
// an error is written and os.exit will be executed.
func Connect(log logger.Logger, conf *config.Config) *gorm.DB {
	log.Infof("Connect to database...")

	db, err := gorm.Open(
		postgres.Open(dsn(conf)),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		},
	)

	if err != nil {
		log.Fatalf("Connect to database DSN: %s failed: %s", dsn(conf), err)
	}

	return db
}

// dsn construct the database dsn from the configuration.
func dsn(config *config.Config) string {
	sslMode := "disable"

	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		config.Db.Test.Host,
		config.Db.Test.Port,
		config.Db.Test.User,
		config.Db.Test.Name,
		config.Db.Test.Password,
		sslMode,
	)
}
