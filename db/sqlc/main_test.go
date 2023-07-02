package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	testDB.SetMaxOpenConns(90)

	if err != nil {
		log.Fatal("cannot create connection to the database", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
