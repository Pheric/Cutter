package main

import (
	"cutter/server"
	"flag"
	"log"
	"cutter/data"
)

const VERSION = 0.01

var createDb bool
var user, pass, db string
var port int

func main() {
	log.Printf("Cutter v%.2f (c) Eric Graham, 2018\n\n", VERSION)
	setupFlags()

	data.InitDb(user, pass, db)
	if createDb {
		data.SetupTables()
	}

	server.Init(port)
}

func setupFlags() {
	flag.BoolVar(&createDb, "createdb", false, "Creates the database")
	flag.StringVar(&db, "database", "cutter", "Database name")
	flag.StringVar(&user, "username", "postgres", "Database username")
	flag.StringVar(&pass, "password", "postgres", "Database password")
	flag.IntVar(&port, "httpport", 2500, "The port the HTTP server listens on")

	flag.Parse()
}
