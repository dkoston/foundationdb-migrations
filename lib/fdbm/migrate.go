package fdbm

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/directory"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
    "github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
)

var (
	NAME = "fdbm"
	ErrNoPreviousVersion = errors.New("no previous version found")
	FdbmDirectory        = "fdbm"			// foundationdb directory where data for fdbm is stored
	MigrationsSubspace   = "migrations"		// subspace where migration data is stored
    STATUS_SUCCESS int64 = 1
    STATUS_FAILED int64 = 2
    STATUS_ROLLEDBACK int64 = 3
)

var migrationsSS subspace.Subspace

type MigrationRecord struct {
	VersionId int64
	TStamp    time.Time
	IsApplied bool // was this a result of up() or down()
}

type Migration struct {
    Name string
    Number int64
    Date string
    Status int64
    StatusName string
    Next int64
    Previous int64
    Source string
}

type migrationSorter []*Migration

// helpers so we can use pkg sort
func (ms migrationSorter) Len() int           { return len(ms) }
func (ms migrationSorter) Swap(i, j int)      { ms[i], ms[j] = ms[j], ms[i] }
func (ms migrationSorter) Less(i, j int) bool { return ms[i].Number < ms[j].Number }

func newMigration(number int64, name string, source string) *Migration {
	return &Migration{Number: number, Name: name, Source: source}
}

func RunMigrations(conf *DBConf, migrationsDir string, targetNumber int64) (err error) {

	db, err := OpenDBFromDBConf(conf)
	if err != nil {
		return err
	}

	return RunMigrationsOnDb(conf, migrationsDir, targetNumber, db)
}

func GetSubspace(db fdb.Database) (migrationsSS subspace.Subspace) {
    fdbmDir, err := directory.CreateOrOpen(db, []string{FdbmDirectory}, nil)
    if err != nil {
        log.Fatal(err)
    }

    return fdbmDir.Sub(MigrationsSubspace)
}


// Runs migration on a specific database instance.
func RunMigrationsOnDb(conf *DBConf, migrationsDir string, targetNumber int64, db fdb.Database) (err error) {
	migrationsSS = GetSubspace(db)

	current, err := GetLatestMigrationFromDB(db, migrationsSS)
	if err != nil {
		return err
	}

	migrations, err := GetMigrationFilesFromDisk(migrationsDir, current.Number, targetNumber)
	if err != nil {
		return err
	}

	if len(migrations) == 0 {
		fmt.Printf("%s: no migrations to run. current version: %d\n", NAME, current)
		return nil
	}

	ms := migrationSorter(migrations)
	direction := current.Number < targetNumber
	ms.Sort(direction)

	fmt.Printf("%s: migrating database. current version: %d, targetNumber: %d\n", NAME, current.Number, targetNumber)

	for _, m := range ms {

		switch filepath.Ext(m.Source) {
		case ".go":
			err = runGoMigration(conf, m.Source, m.Number, direction)
		}

		if err != nil {
			return errors.New(fmt.Sprintf("FAIL %v, quitting migration", err))
		}

		fmt.Println("OK   ", filepath.Base(m.Source))
	}

	return nil
}

// collect all the valid looking migration scripts in the
// migrations folder, and key them by version
func GetMigrationFilesFromDisk(dirpath string, current, target int64) (m []*Migration, err error) {

	// extract the numeric component of each migration,
	// filter out any uninteresting files,
	// and ensure we only have one file per migration version.
	filepath.Walk(dirpath, func(name string, info os.FileInfo, err error) error {

		if number, e := NumericComponent(name); e == nil {
			for _, g := range m {
				if number == g.Number {
					log.Fatalf("more than one file specifies the migration for version %d (%s and %s)",
                        number, g.Source, name)
				}
			}

			if versionFilter(number, current, target) {
				m = append(m, newMigration(number, name,  name))
			}
		}

		return nil
	})

	return m, nil
}

func versionFilter(v, current, target int64) bool {

	if target > current {
		return v > current && v <= target
	}

	if target < current {
		return v <= current && v > target
	}

	return false
}

func (ms migrationSorter) Sort(direction bool) {

	// sort ascending or descending by version
	if direction {
		sort.Sort(ms)
	} else {
		sort.Sort(sort.Reverse(ms))
	}

	// now that we're sorted in the appropriate direction,
	// populate next and previous for each migration
	for i, m := range ms {
		prev := int64(-1)
		if i > 0 {
			prev = ms[i-1].Number
			ms[i-1].Next = m.Number
		}
		ms[i].Previous = prev
	}
}

// look for migration scripts with names in the form:
//  XXX_descriptivename.ext
// where XXX specifies the version number
// and ext specifies the type of migration
func NumericComponent(name string) (int64, error) {

	base := filepath.Base(name)

	if ext := filepath.Ext(base); ext != ".go" {
		return 0, errors.New("not a recognized migration file type")
	}

	idx := strings.Index(base, "_")
	if idx < 0 {
		return 0, errors.New("no separator found")
	}

	n, e := strconv.ParseInt(base[:idx], 10, 64)
	if e == nil && n <= 0 {
		return 0, errors.New("migration IDs must be greater than zero")
	}

	return n, e
}

// retrieve the last applied migration name
func GetLatestMigrationFromDB(t fdb.Transactor, migrationsSS subspace.Subspace) (latest Migration, err error) {
    latest = Migration{Number: 000}

    migrations, err := GetMigrationsFromDB(t, migrationsSS)

    if err != nil {
        log.Printf("Unable to load migrations from database: %v", err)
        return latest, err
    }

    for _, migration := range migrations {
        if migration.Number > latest.Number {
            if migration.Status == STATUS_SUCCESS {
                latest = migration
            }
        }
    }

	return
}

func mapMigrationStatusToStatusName(statusName string) string {
    num, err := strconv.ParseInt(statusName, 10, 64)
    if err != nil {
        log.Fatalf("Unable to parse migration status, %v", err)
    }
    switch num {
    case STATUS_SUCCESS:
        return "Success"
    case STATUS_FAILED:
        return "Failed"
    case STATUS_ROLLEDBACK:
        return "Rolled back"
    default:
        return "Unknown"
    }
}

func parseMigrationStatus(status string) int64 {
    num, err := strconv.ParseInt(status, 10, 64)
    if err != nil {
        log.Fatalf("Unable to parse migration status, %v", err)
    }
    return num
}

func GetMigrationsFromDB(t fdb.Transactor, migrationsSS subspace.Subspace) (migrations map[int64]Migration, err error) {
    migrations = make(map[int64]Migration)
    r, err := t.ReadTransact(func (rtr fdb.ReadTransaction) (interface{}, error) {
        ri := rtr.GetRange(migrationsSS, fdb.RangeOptions{}).Iterator()
        for ri.Advance() {
            kv := ri.MustGet()
            k, err := migrationsSS.Unpack(kv.Key)
            if err != nil {
                log.Fatalf("Unable to parse KV, %v", err)
            }
            i, ok := k[0].(int64)

            if !ok {
                log.Printf("Unable to convert %v to int64 %v\n", k[0], i)
                continue
            }
            migration, isset := migrations[i]

            if !isset {
                migration = Migration{}
            }

            v := string(kv.Value)
            switch key := k[1]; key {
            case "Date":
                migration.Date = v
            case "Name":
                migration.Name = v
                migration.Number, err = NumericComponent(v)
            case "Status":
                migration.Status = parseMigrationStatus(v)
                migration.StatusName = mapMigrationStatusToStatusName(v)
            }
            migrations[i] = migration
        }

        return migrations, nil
    })

    if err == nil {
        migrations = r.(map[int64]Migration)
    }

    return migrations, err
}

func FilterMigrationsByNumber(migrations map[int64]Migration, number int64) (migration Migration) {
    migration = Migration{Number: 000}

    for _, m := range migrations {
        if m.Number == number {
            return m
        }
    }
    return migration
}


// Gets the version of the latest migration that was successful
func GetDBVersion(t fdb.Transactor, migrationsSS subspace.Subspace) (version int64, err error) {
	migration, err := GetLatestMigrationFromDB(t, migrationsSS)
	if err != nil {
		return -1, err
	}

	return migration.Number, nil
}

func GetPreviousDBVersion(dirpath string, version int64) (previous int64, err error) {

	previous = -1
	sawGivenVersion := false

	filepath.Walk(dirpath, func(name string, info os.FileInfo, walkerr error) error {

		if !info.IsDir() {
			if v, e := NumericComponent(name); e == nil {
				if v > previous && v < version {
					previous = v
				}
				if v == version {
					sawGivenVersion = true
				}
			}
		}

		return nil
	})

	if previous == -1 {
		if sawGivenVersion {
			// the given version is (likely) valid but we didn't find
			// anything before it.
			// 'previous' must reflect that no migrations have been applied.
			previous = 0
		} else {
			err = ErrNoPreviousVersion
		}
	}

	return
}

// helper to identify the most recent possible version
// within a folder of migration scripts
func GetMostRecentDBVersion(dirpath string) (version int64, err error) {

	version = -1

	filepath.Walk(dirpath, func(name string, info os.FileInfo, walkerr error) error {
		if walkerr != nil {
			return walkerr
		}

		if !info.IsDir() {
			if v, e := NumericComponent(name); e == nil {
				if v > version {
					version = v
				}
			}
		}

		return nil
	})

	if version == -1 {
		err = errors.New("no valid version found")
	}

	return
}

func CreateMigration(name, dir string, t time.Time) (path string, err error) {

	timestamp := FormatTime(t)
	filename := fmt.Sprintf("%v_%v.%v", timestamp, name, "go")

	fpath := filepath.Join(dir, filename)
	path, err = writeTemplateToFile(fpath, goMigrationTemplate, timestamp)

	return
}


func FormatTime(t time.Time) (date string) {
    return t.Format("20060102150405")
}


func FinalizeMigration(t fdb.Transactor, migrationsSS subspace.Subspace, date string, name string, number int64, statusInt int64) (err error) {
    _, err = t.Transact(func (tr fdb.Transaction) (interface{}, error) {
        status := strconv.FormatInt(statusInt, 10)

        tr.Set(migrationsSS.Pack(tuple.Tuple{number, "Name"}), []byte(name))
        tr.Set(migrationsSS.Pack(tuple.Tuple{number, "Number"}), []byte(string(number)))
        tr.Set(migrationsSS.Pack(tuple.Tuple{number, "Date"}), []byte(date))
        tr.Set(migrationsSS.Pack(tuple.Tuple{number, "Status"}), []byte(status))
        tr.Set(migrationsSS.Pack(tuple.Tuple{number, "StatusName"}), []byte(mapMigrationStatusToStatusName(status)))

        return nil, nil
    })

    return err
}


var goMigrationTemplate = template.Must(template.New("goose.go-migration").Parse(`
package main

import (
	"github.com/apple/foundationdb/bindings/go/src/fdb"
)

// Up is executed when this migration is applied
func Up_{{ . }}(t fdb.Transactor) error {
	_, err := t.Transact(func (tr fdb.Transaction) (ret interface{}, err error) {
		// Add your code here, i.e. tr.Set(fdb.Key("hello"), []byte("world"))
        return
    })

    return err
}

// Down is executed when this migration is rolled back
func Down_{{ . }}(t fdb.Transactor) error {
	_, err := t.Transact(func (tr fdb.Transaction) (ret interface{}, err error) {
		// Add your code here. i.e. tr.Clear(fdb.Key("hello"))
        return
    })

    return err
}
`))