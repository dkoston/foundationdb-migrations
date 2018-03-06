package main

import (
	"github.com/cryptowalkio/goose/lib/goose"
	"fmt"
	"log"
)

var dbVersionCmd = &Command{
	Name:    "dbversion",
	Usage:   "dbversion [migration_table_prefix]",
	Summary: "Print the current version of the database",
	Help:    `dbversion extended help here...`,
	Run:     dbVersionRun,
}

func dbVersionRun(cmd *Command, args ...string) {
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

	fmt.Printf("goose: dbversion %v\n", current)
}
