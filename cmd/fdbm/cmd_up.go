package main

import (
	"github.com/dkoston/foundationdb-migrations/lib/fdbm"
	"log"
)

var upCmd = &Command{
	Name:    "up",
	Usage:   "up [migration_table_prefix]",
	Summary: "Migrate the DB to the most recent version available",
	Help:    `migration_table_prefix should be ^[a-z_]+$`,
	Run:     upRun,
}

func upRun(cmd *Command, args ...string) {
	conf, err := dbConfFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	target, err := fdbm.GetMostRecentDBVersion(conf.MigrationsDir)
	if err != nil {
		log.Fatal(err)
	}

	if err := fdbm.RunMigrations(conf, conf.MigrationsDir, target); err != nil {
		log.Fatal(err)
	}
}
