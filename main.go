package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/johnquangdev/oauth2/cmd"
	"github.com/johnquangdev/oauth2/cmd/sqlmigrate"
)

func main() {
	// Define CLI flags
	migrateDown := flag.Bool("migrate-down", false, "Run database migration rollback")
	migrateLimit := flag.Int("limit", 1, "Number of migrations to rollback (default: 1)")
	flag.Parse()

	// Check if migrate-down flag is set
	if *migrateDown {
		fmt.Printf("Running migration down (limit: %d)...\n", *migrateLimit)
		sqlmigrate.RunMigrateDown(*migrateLimit)
		os.Exit(0)
	}

	// Run normal server
	cmd.Run()
}
