package main

import (
	"github.com/dkoston/foundationdb-migrations/lib/fdbm"
	"fmt"
	"log"
)
var NAME = "fdbm"

var statusCmd = &Command{
	Name:    "status",
	Usage:   "",
	Summary: "Check the status of the database versus migration files on disk",
	Help:    `status extended help here...`,
	Run:     statusRun,
}

func statusRun(cmd *Command, args ...string) {
	conf, err := dbConfFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	// collect all migrationFiles
	min := int64(0)
	max := int64((1 << 63) - 1)
	migrationFiles, e := fdbm.GetMigrationFilesFromDisk(conf.MigrationsDir, min, max)
	if e != nil {
		log.Fatal(e)
	}

	db, e := fdbm.OpenDBFromDBConf(conf)
	if e != nil {
		log.Fatal("couldn't open DB:", e)
	}

    migrationsSS := fdbm.GetSubspace(db)

	// check if we have any migrations in the database
	dbMigrations, e := fdbm.GetMigrationsFromDB(db, migrationsSS)
	if e != nil {
		log.Fatal(e)
	}

	fmt.Printf("%s: migration status\n", NAME)
	fmt.Println("    Status         Date                      File                            ")
	fmt.Println("    =========================================================================")
	for _, m := range migrationFiles {
		printMigrationStatus(dbMigrations, m)
	}
}

func printMigrationStatus(dbMigrations map[int64]fdbm.Migration, fileMigration *fdbm.Migration) {

    dbMigration := fdbm.FilterMigrationsByNumber(dbMigrations, fileMigration.Number)

    var status = "Pending"
    var date = "n/a"
    var fileName = fileMigration.Name

    if dbMigration.Number != 000 {
        status = dbMigration.StatusName
        date = dbMigration.Date
    }

	fmt.Printf("    %s        %-24s %v\n", status, date, fileName)
}