package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_DSN"))
	if err != nil {
		log.Fatalln(err)
	}

	_, err = FindUsers(db)
	if err != nil {
		log.Fatalln(err)
	}
}
