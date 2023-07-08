package main

import (
	"database/sql"
	"github.com/HBeserra/simplebank/api"
	db "github.com/HBeserra/simplebank/db/sqlc"
	_ "github.com/lib/pq"
	"log"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func main() {

	// Todo: add .env configuration load

	conn, err := sql.Open(dbDriver, dbSource)
	conn.SetMaxOpenConns(90)

	if err != nil {
		log.Fatal("cannot create connection to the database", err)
	}

	store := db.NewStore(conn)

	server := api.NewServer(store)
	err = server.Start("0.0.0.0:3306")

	if err != nil {
		log.Fatal("Cannot start server:", err.Error())
	}
}
