package main

import (
	"database/sql"
	"log"

	"github.com/DamirAbd/HLA-HW1/cmd/api"
	"github.com/DamirAbd/HLA-HW1/db"
)

func main() {
	connStr := "postgres://root:pwd@pgres:5432/socntw?sslmode=disable"
	tdb, err := db.SQLStorage(connStr)
	if err != nil {
		log.Fatal(err)
	}

	connStrC := "postgres://postgres:pwd@citus_coordinator:5432/cdb?sslmode=disable"
	cdb, err := db.SQLStorage(connStrC)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(tdb)
	initStorage(cdb)

	server := api.NewAPIServer(":8080", tdb, cdb)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}
