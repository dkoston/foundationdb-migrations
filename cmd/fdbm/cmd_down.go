package main

import (
	"github.com/dkoston/foundationdb-migrations/lib/fdbm"
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

	previous, err := fdbm.GetPreviousDBVersion(conf.MigrationsDir, current)
	if err != nil {
		log.Fatal(err)
	}

	if err = fdbm.RunMigrations(conf, conf.MigrationsDir, previous); err != nil {
		log.Fatal(err)
	}
}
