package fdbm

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
)

// DBDriver encapsulates the info needed to work with
// a specific database driver
type DBDriver struct {
	Name    string
	OpenStr string
	Import  string
}

type DBConf struct {
	MigrationsDir string
	ClusterFile   string
	FDBAPIVersion int
}

func dirOrFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
	if os.IsNotExist(err) { return false, nil }
	return true, err
}


func NewDBConf(migrationsDir, clusterFile string, fdbVersion int) (*DBConf, error) {

	// Make sure the migrations directory exists
	exists, err := dirOrFileExists(migrationsDir)

	if !exists {
		errString := fmt.Sprintf("Migrations directory does not exist: %s", migrationsDir)
		return nil, errors.New(errString)
	}

	if err != nil {
		return nil, err
	}

	// Make sure the cluster file exists
	cFileExists, err := dirOrFileExists(clusterFile)

	if !cFileExists {
		errString := fmt.Sprintf("FoundationDB cluster file does not exist: %s", clusterFile)
		return nil, errors.New(errString)
	}

	if err != nil {
		return nil, err
	}

	return &DBConf{
		MigrationsDir: migrationsDir,
		ClusterFile: clusterFile,
		FDBAPIVersion: fdbVersion,
	}, nil
}


func OpenDBFromDBConf(conf *DBConf) (fdb.Database, error) {
	fdb.MustAPIVersion(conf.FDBAPIVersion)

	// Open the default database from the system cluster
	db, err := fdb.Open(conf.ClusterFile, []byte("DB"))

	if err != nil {
		log.Fatalf("Unable to open database: %v", err)
	}

	return db, nil
}
