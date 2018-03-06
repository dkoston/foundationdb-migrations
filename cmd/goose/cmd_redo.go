package main

import (
	"github.com/cryptowalkio/goose/lib/goose"
	"log"
)

var redoCmd = &Command{
	Name:    "redo",
	Usage:   "redo [migration_table_prefix]",
	Summary: "Re-run the latest migration",
	Help:    `migration_table_prefix should be ^[a-z_]+$`,
	Run:     redoRun,
}

func redoRun(cmd *Command, args ...string) {
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

	if err := goose.RunMigrations(conf, conf.MigrationsDir, previous, tablePrefix); err != nil {
		log.Fatal(err)
	}

	if err := goose.RunMigrations(conf, conf.MigrationsDir, current, tablePrefix); err != nil {
		log.Fatal(err)
	}
}
