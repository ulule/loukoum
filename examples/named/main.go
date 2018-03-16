package main

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatalln(err)
	}

	_, err = FindUsers(db)
	if err != nil {
		log.Fatalln(err)
	}
}
