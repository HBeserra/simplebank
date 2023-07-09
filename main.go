package main

import (
	"database/sql"
	"github.com/HBeserra/simplebank/api"
	db "github.com/HBeserra/simplebank/db/sqlc"
	"github.com/HBeserra/simplebank/util"
	_ "github.com/lib/pq"
	"log"
)

func main() {

	// Todo: add .env configuration load

	config, err := util.LoadConfig(".")

	conn, err := sql.Open(config.DBDriver, config.DBUri)
	conn.SetMaxOpenConns(config.DBMaxNConn)

	if err != nil {
		log.Fatal("cannot create connection to the database", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("Cannot start server:", err.Error())
	}
}
