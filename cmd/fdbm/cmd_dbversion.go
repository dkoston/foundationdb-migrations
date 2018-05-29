package main

import (
	"github.com/dkoston/foundationdb-migrations/lib/fdbm"
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
	conf, err := dbConfFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	db, e := fdbm.OpenDBFromDBConf(conf)
	if e != nil {
		log.Fatal("couldn't open DB:", e)
	}

	migrationsSS := fdbm.GetSubspace(db)

	current, err := fdbm.GetDBVersion(db, migrationsSS)
	if err != nil {
		log.Fatal(err)
	}

	if current == 000 {
		log.Fatal("no migrations applied. Use 'goose up' to apply migration files")
	}

	fmt.Printf("goose: dbversion %v\n", current)
}
