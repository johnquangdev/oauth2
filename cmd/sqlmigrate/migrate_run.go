package sqlmigrate

import (
	"fmt"
	"log"

	"github.com/johnquangdev/oauth2/utils"
	migrate "github.com/rubenv/sql-migrate"
	"gorm.io/gorm"
)

func RunSqlMigrate(cfg utils.Config, db *gorm.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "sqlmigrations",
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB, error: %w", err)
	}

	total, err := migrate.Exec(sqlDB, "postgres", migrations, migrate.Up)
	if err != nil {
		return fmt.Errorf("cannot execute migration: %w", err)
	}

	log.Printf("applied %d migrations\n", total)

	return nil
}

func DownSqlMigrate(cfg utils.Config, db *gorm.DB, limit int) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "sqlmigrations",
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB, error: %w", err)
	}

	total, err := migrate.ExecMax(sqlDB, "postgres", migrations, migrate.Down, limit)
	if err != nil {
		return fmt.Errorf("cannot down migration: %w", err)
	}

	log.Printf("applied %d migrations\n", total)

	return nil
}
