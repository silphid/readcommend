package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
	"github.com/kelseyhightower/envconfig"
)

const usageText = `This program runs command on the db. Supported commands are:
  - init - creates version info table in the database
  - up - runs all available migrations.
  - up [target] - runs available migrations up to the target one.
  - down - reverts last migration.
  - reset - reverts all migrations.
  - version - prints current db version.
  - set_version [version] - sets db version without running migrations.

Usage:
  go run migrations/main.go <command> [args]
`

type config struct {
	DbURL string `envconfig:"DB_URL" required:"true"`
}

// main leverages go-pg/migrations to apply various database migration operations
func main() {
	// Parse command-line
	flag.Usage = func() {
		fmt.Print(usageText)
		os.Exit(0)
	}
	flag.Parse()
	subCmd := flag.Arg(0)

	// Connect to database
	cfg := loadConfig()
	dbOptions, err := pg.ParseURL(cfg.DbURL)
	if err != nil {
		log.Fatalln(err)
	}
	db := pg.Connect(dbOptions)

	// Lock across all replicas to prevent concurrent migrations
	unlock := getAdvisoryLock(db)
	defer unlock()

	// Apply migrations
	migrations.SetTableName("migrations")
	oldVersion, newVersion, err := migrations.Run(db, flag.Args()...)
	if err != nil {
		panic(err)
	}

	// The "init" sub-command always returns old and new versions as 0 so we
	// need to run "version" to get actual database version for logging purposes
	if subCmd == "init" {
		oldVersion, newVersion, err = migrations.Run(db, "version")
		if err != nil {
			panic(err)
		}
	}

	// Log performed migrations, if any
	if newVersion != oldVersion {
		fmt.Printf("%s: migrated from version %d to %d\n", subCmd, oldVersion, newVersion)
	} else {
		fmt.Printf("%s: version is %d\n", subCmd, oldVersion)
	}
}

// getAdvisoryLock prevents multiple replicas running against same database from concurrently applying migrations.
// It returns a function to call to release the lock.
func getAdvisoryLock(database *pg.DB) func() {
	const lockID = 8913520
	_, err := database.Exec(fmt.Sprintf("select pg_advisory_lock(%d);", lockID))
	if err != nil {
		panic(err)
	}

	return func() {
		_, err := database.Exec(fmt.Sprintf("select pg_advisory_unlock(%d);", lockID))
		if err != nil {
			panic(err)
		}
	}
}

// loadConfig loads the... (you have 3 seconds to guess...) config!
func loadConfig() config {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}
