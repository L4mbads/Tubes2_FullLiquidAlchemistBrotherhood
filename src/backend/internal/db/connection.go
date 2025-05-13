package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	// Configure the connection pool
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(5)
	// db.SetConnMaxLifetime(0)

	return db
}

func ClearDB(dbConn *sql.DB) {
	// create table if it doesn't exist
	_, err := dbConn.Exec("DROP TABLE IF EXISTS recipes")
	if err != nil {
		log.Fatal(err)
	}
	_, err = dbConn.Exec("DROP TABLE IF EXISTS elements")
	if err != nil {
		log.Fatal(err)
	}
	_, err = dbConn.Exec("CREATE TABLE IF NOT EXISTS elements (name TEXT, image_url TEXT, type SMALLINT)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = dbConn.Exec("CREATE TABLE IF NOT EXISTS recipes (element TEXT, ingredient1 TEXT, ingredient2 TEXT)")
	if err != nil {
		log.Fatal(err)
	}
}
