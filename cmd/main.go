package main

import (
	"database/sql"
	"log"

	"github.com/DamirAbd/HLA-HW1/cmd/api"
	"github.com/DamirAbd/HLA-HW1/db"
)

func main() {
	connStr := "postgres://root:pwd@pgres:5432/socntw?sslmode=disable"
	db, err := db.SQLStorage(connStr)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(":8080", db)
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
