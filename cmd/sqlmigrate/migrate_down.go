package sqlmigrate

import (
	"log"

	"github.com/johnquangdev/oauth2/repository"
	"github.com/johnquangdev/oauth2/utils"
)

// RunMigrateDown chạy migration rollback
func RunMigrateDown(limit int) {
	// Load config
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect database
	db, err := repository.ConnectPostgres(*config)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Run migrate down
	err = DownSqlMigrate(*config, db, limit)
	if err != nil {
		log.Fatalf("Failed to run SQL migration down: %v", err)
	}

	log.Println("Migration down completed successfully!")
}
