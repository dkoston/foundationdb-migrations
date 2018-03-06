package main

import (
	"github.com/cryptowalkio/goose/lib/goose"
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
	tablePrefix := "goose"

	if len(args) >= 1 {
		tablePrefix = args[0]
	}

	conf, err := dbConfFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	target, err := goose.GetMostRecentDBVersion(conf.MigrationsDir)
	if err != nil {
		log.Fatal(err)
	}

	if err := goose.RunMigrations(conf, conf.MigrationsDir, target, tablePrefix); err != nil {
		log.Fatal(err)
	}
}
