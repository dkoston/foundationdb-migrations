package main

import (
	"github.com/cryptowalkio/goose/lib/goose"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"time"
	"bytes"
)

var statusCmd = &Command{
	Name:    "status",
	Usage:   "",
	Summary: "dump the migration status for the current DB",
	Help:    `status extended help here...`,
	Run:     statusRun,
}

type StatusData struct {
	Source string
	Status string
}

func statusRun(cmd *Command, args ...string) {
	tablePrefix := "goose"

	if len(args) >= 1 {
		tablePrefix = args[0]
	}

	conf, err := dbConfFromFlags()
	if err != nil {
		log.Fatal(err)
	}

	// collect all migrations
	min := int64(0)
	max := int64((1 << 63) - 1)
	migrations, e := goose.CollectMigrations(conf.MigrationsDir, min, max)
	if e != nil {
		log.Fatal(e)
	}

	db, e := goose.OpenDBFromDBConf(conf)
	if e != nil {
		log.Fatal("couldn't open DB:", e)
	}
	defer db.Close()

	// must ensure that the version table exists if we're running on a pristine DB
	if _, e := goose.EnsureDBVersion(conf, db, tablePrefix); e != nil {
		log.Fatal(e)
	}

	fmt.Printf("goose: status for environment '%v'\n", conf.Env)
	fmt.Println("    Applied At                  Migration")
	fmt.Println("    =======================================")
	for _, m := range migrations {
		printMigrationStatus(db, m.Version, filepath.Base(m.Source), tablePrefix)
	}
}

func printMigrationStatus(db *sql.DB, version int64, script string, tablePrefix string) {
	var row goose.MigrationRecord
	tableName := getVersionTableName(tablePrefix)
	statement := fmt.Sprintf("SELECT tstamp, is_applied FROM %s WHERE version_id=%d ORDER BY tstamp DESC LIMIT 1", tableName, version)
	e := db.QueryRow(statement).Scan(&row.TStamp, &row.IsApplied)

	if e != nil && e != sql.ErrNoRows {
		log.Fatal(e)
	}

	var appliedAt string

	if row.IsApplied {
		appliedAt = row.TStamp.Format(time.ANSIC)
	} else {
		appliedAt = "Pending"
	}

	fmt.Printf("    %-24s -- %v\n", appliedAt, script)
}

// get the tablename used for migrations
func getVersionTableName(tablePrefix string) string {
	var buffer bytes.Buffer
	buffer.WriteString(tablePrefix)
	buffer.WriteString("_db_version")
	return buffer.String()
}