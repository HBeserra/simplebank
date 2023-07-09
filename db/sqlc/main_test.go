package db

import (
	"database/sql"
	"github.com/HBeserra/simplebank/util"
	_ "github.com/lib/pq"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatalln("can't load the configuration", err.Error())
	}

	testDB, err = sql.Open(config.DBDriver, config.DBUri)
	testDB.SetMaxOpenConns(config.DBMaxNConn)

	if err != nil {
		log.Fatal("cannot create connection to the database", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())
}
