package data

import (
	"fmt"
	"github.com/go-pg/pg"
	"log"
	"time"
)

var conn *pg.Options

func InitDb(username, password, database string) {
	conn = &pg.Options{
		User:        username,
		Password:    password,
		Database:    database,
		MaxRetries:  3,
		PoolSize:    5,
		PoolTimeout: 4 * time.Second,
		IdleTimeout: 15 * time.Minute,
	}
}

func openConnection() *pg.DB {
	return pg.Connect(conn)
}

func SetupTables() {
	log.Println("Creating tables...")
	start := time.Now()
	conn := openConnection()

	tables := map[string]string{
		// Ignore column Employees.payments
		"Employees": `
			uuid UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
			first VARCHAR DEFAULT 'John',
			last VARCHAR DEFAULT 'Doe',
			phone VARCHAR DEFAULT '+1 (000) 000-0000',
			owed NUMERIC(6, 2) DEFAULT 0
		`,
		// Ignore columns clients.payments and clients.cuts
		"clients": `
			uuid UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
			first VARCHAR DEFAULT 'John',
			last VARCHAR DEFAULT 'Doe',
			phone VARCHAR DEFAULT '+1 (000) 000-0000',
			address VARCHAR DEFAULT '',
			quote NUMERIC(6, 2) DEFAULT 30,
			ttc SMALLINT DEFAULT 15,
			period SMALLINT DEFAULT 0,
			balance NUMERIC(6, 2) DEFAULT 0
		`,
		"cuts": `
			uuid UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
			client UUID,
			date DATE DEFAULT 'today',
			price NUMERIC(6, 2) DEFAULT 30,
			Employees UUID[]
		`,
		"payments": `
			uuid UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
			target UUID,
			amount NUMERIC(6, 2),
			date DATE DEFAULT 'today'
		`,
	}

	// Create extension for uuid generation
	_, err := conn.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
	if err != nil {
		log.Printf("Error creating extension 'uuid-ossp' (for uuid generation): %v\n", err)
	}

	// Not using the ORM's solution because it looks super messy and hacky.
	errCount := 0
	for l, t := range tables {
		if _, err := conn.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", l, t)); err != nil {
			log.Printf("Error creating table %s: %v\n", l, err)
			errCount++
		}
	}

	log.Printf("Finished creating %d tables in %gms\n", len(tables)-errCount, time.Since(start).Seconds()/1000)
}

func LogDbQueries(db *pg.DB) {
	db.OnQueryProcessed(func(event *pg.QueryProcessedEvent) {
		query, err := event.FormattedQuery()
		if err != nil {
			panic(err)
		}

		log.Printf("%s %s", time.Since(event.StartTime), query)
	})
}
