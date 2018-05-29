package main

import (
	"github.com/dkoston/foundationdb-migrations/lib/fdbm"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
)

// global options. available to any subcommands.
var migrationsDir = flag.String("m", "db/migrations", "folder containing your migrations (default = ./db/migrations)")
var clusterFile = flag.String("f", "db/fdb.cluster", "path to your FoundationDB cluster file (default = ./db/fdb.cluster)")
var fdbVersion = flag.Int("v", 510, "Version of the FoundationDB Cluster API (default = 510)")

// helper to create a DBConf from the given flags
func dbConfFromFlags() (dbconf *fdbm.DBConf, err error) {
	return fdbm.NewDBConf(*migrationsDir, *clusterFile, *fdbVersion)
}

var commands = []*Command{
	upCmd,
	downCmd,
	redoCmd,
	statusCmd,
	createCmd,
	dbVersionCmd,
}

func main() {

	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 || args[0] == "-h" {
		flag.Usage()
		return
	}

	var cmd *Command
	name := args[0]
	for _, c := range commands {
		if strings.HasPrefix(c.Name, name) {
			cmd = c
			break
		}
	}

	if cmd == nil {
		fmt.Printf("error: unknown command %q\n", name)
		flag.Usage()
		os.Exit(1)
	}

	cmd.Exec(args[1:])
}

func usage() {
	fmt.Print(usagePrefix)
	flag.PrintDefaults()
	usageTmpl.Execute(os.Stdout, commands)
}

var usagePrefix = `
fdbm is a database migration tool for FoundationDB (https://github.com/apple/foundationdb).

Usage:
    fdbm [options] <subcommand> [subcommand options]

Options:
`
var usageTmpl = template.Must(template.New("usage").Parse(
	`
Commands:{{range .}}
    {{.Name | printf "%-10s"}} {{.Summary}}{{end}}
`))
