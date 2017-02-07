package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/smacker/goose"

	_ "github.com/mattn/go-sqlite3"
)

var (
	flags  = flag.NewFlagSet("goose", flag.ExitOnError)
	dir    = flags.String("dir", ".", "directory with migration files")
	driver = "sqlite3"
)

func main() {
	flags.Usage = usage
	flags.Parse(os.Args[1:])

	args := flags.Args()

	if len(args) > 1 && args[0] == "create" {
		if err := goose.Run("create", nil, *dir, args[1:]...); err != nil {
			log.Fatalf("goose run: %v", err)
		}
		return
	}

	if len(args) < 2 {
		flags.Usage()
		return
	}

	if args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	dbstring, command := args[0], args[1]

	if err := goose.SetDialect(driver); err != nil {
		log.Fatal(err)
	}

	switch dbstring {
	case "":
		log.Fatalf("-dbstring=%q not supported\n", dbstring)
	default:
	}

	db, err := sql.Open(driver, dbstring)
	if err != nil {
		log.Fatalf("-dbstring=%q: %v\n", dbstring, err)
	}

	arguments := []string{}
	if len(args) > 2 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose run: %v", err)
	}
}

func usage() {
	fmt.Print(usagePrefix)
	flags.PrintDefaults()
	fmt.Print(usageCommands)
}

var (
	usagePrefix = `Usage: goose [OPTIONS] DRIVER DBSTRING COMMAND

Examples:
    goose ./foo.db status
    goose create init sql

Options:
`

	usageCommands = `
Commands:
	apply      Apply all pending migrations
	reset      Rollback all database migrations
	refresh    Reset and re-run all migrations
	up         Migrate the DB to the most recent version available
	down       Roll back the version by 1
	redo       Re-run the latest migration
	status     Dump the migration status for the current DB
	dbversion  Print the current version of the database
	create     Creates a blank migration template
`
)
