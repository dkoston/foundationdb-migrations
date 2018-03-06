package main

import (
	"github.com/cryptowalkio/goose/lib/goose"
	"log"
)

var downCmd = &Command{
	Name:    "down",
	Usage:   "down [migration_table_prefix]",
	Summary: "Roll back the version by 1",
	Help:    `migration_table_prefix should be ^[a-z_]+$`,
	Run:     downRun,
}

func downRun(cmd *Command, args ...string) {
	tablePrefix := "goose"

	if len(args) >= 1 {
		tablePrefix = args[0]
	}

	conf, err := dbConfFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	current, err := goose.GetDBVersion(conf, tablePrefix)
	if err != nil {
		log.Fatal(err)
	}

	previous, err := goose.GetPreviousDBVersion(conf.MigrationsDir, current)
	if err != nil {
		log.Fatal(err)
	}

	if err = goose.RunMigrations(conf, conf.MigrationsDir, previous, tablePrefix); err != nil {
		log.Fatal(err)
	}
}
