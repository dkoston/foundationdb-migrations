package main

import (
	"github.com/dkoston/foundationdb-migrations/lib/fdbm"
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

	if err := fdbm.RunMigrations(conf, conf.MigrationsDir, previous); err != nil {
		log.Fatal(err)
	}

	if err := fdbm.RunMigrations(conf, conf.MigrationsDir, current); err != nil {
		log.Fatal(err)
	}
}
