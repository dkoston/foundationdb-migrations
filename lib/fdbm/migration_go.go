package fdbm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"encoding/gob"
)

type templateData struct {
	Number     int64
	Name       string
	Conf       string // gob encoded DBConf
	Status     int64
	Func       string
}

//
// Run a .go migration.
//
// In order to do this, we copy a modified version of the
// original .go migration, and execute it via `go run` along
// with a main() of our own creation.
//
func runGoMigration(conf *DBConf, path string, number int64, direction bool) error {

	// everything gets written to a temp dir, and zapped afterwards
	d, e := ioutil.TempDir("", "fdbm")
	if e != nil {
		log.Fatal(e)
	}
	defer os.RemoveAll(d)

	directionStr := "Down"
	if direction {
		directionStr = "Up"
	}

	status := STATUS_SUCCESS

	if !direction {
	    status = STATUS_ROLLEDBACK
    }


	var bb bytes.Buffer
	if err := gob.NewEncoder(&bb).Encode(conf); err != nil {
		return err
	}

	// XXX: there must be a better way of making this byte array
	// available to the generated code...
	// but for now, print an array literal of the gob bytes
	var sb bytes.Buffer
	sb.WriteString("[]byte{ ")
	for _, b := range bb.Bytes() {
		sb.WriteString(fmt.Sprintf("0x%02x, ", b))
	}
	sb.WriteString("}")

	td := &templateData{
		Number:    	number,
		Name:       path,
		Conf:       sb.String(),
		Status:     status,
		Func:       fmt.Sprintf("%v_%v", directionStr, number),
	}
	main, e := writeTemplateToFile(filepath.Join(d, "fdbm_main.go"), goMigrationDriverTemplate, td)
	if e != nil {
		log.Fatal(e)
	}

	outpath := filepath.Join(d, filepath.Base(path))
	if _, e = copyFile(outpath, path); e != nil {
		log.Fatal(e)
	}

	cmd := exec.Command("go", "run", main, outpath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if e = cmd.Run(); e != nil {
		log.Fatal("`go run` failed: ", e)
	}

	return nil
}

//
// template for the main entry point to a go-based migration.
// this gets linked against the substituted versions of the user-supplied
// scripts in order to execute a migration via `go run`
//
var goMigrationDriverTemplate = template.Must(template.New("goose.go-driver").Parse(`
package main

import (
	"log"
	"bytes"
	"encoding/gob"
    "time"

	"github.com/dkoston/foundationdb-migrations/lib/fdbm"
)

func main() {

	var conf fdbm.DBConf
	buf := bytes.NewBuffer({{ .Conf }})
	if err := gob.NewDecoder(buf).Decode(&conf); err != nil {
		log.Fatal("gob.Decode - ", err)
	}

	db, err := fdbm.OpenDBFromDBConf(&conf)
	if err != nil {
		log.Fatal("failed to open DB:", err)
	}

    migrationsSS := fdbm.GetSubspace(db)

	err = {{ .Func }}(db)
    if err != nil {
        log.Fatal("Migration Failed: %v", err)
    }

    date := fdbm.FormatTime(time.Now())
    name := "{{ .Name }}"
    status := int64({{ .Status }})

    err = fdbm.FinalizeMigration(db, migrationsSS, date, name, {{ .Number }}, status)
    if err != nil {
        log.Fatal("Failed to record version in database:", err)
    }
}
`))
