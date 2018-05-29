package main

import (
	"github.com/dkoston/foundationdb-migrations/lib/fdbm"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var createCmd = &Command{
	Name:    "create",
	Usage:   "create <name>",
	Summary: "Create the scaffolding for a new migration",
	Help:    `Date based versions will be added automatically, i.e. 20180711_<name>`,
	Run:     createRun,
}

func createRun(cmd *Command, args ...string) {

	if len(args) < 1 {
		log.Fatal("fdbm create: migration name required")
	}

	conf, err := dbConfFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	if err = os.MkdirAll(conf.MigrationsDir, 0777); err != nil {
		log.Fatal(err)
	}

	n, err := fdbm.CreateMigration(args[0], conf.MigrationsDir, time.Now())
	if err != nil {
		log.Fatal(err)
	}

	a, e := filepath.Abs(n)
	if e != nil {
		log.Fatal(e)
	}

	fmt.Println("fdbm create: migration created", a)
}
