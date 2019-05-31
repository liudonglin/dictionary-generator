package db

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestPingDatabase(t *testing.T) {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = pingDatabase(db)
	if err != nil {
		log.Fatal(err)
	}
}

func TestSetupDatabase(t *testing.T) {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = setupDatabase(db, "mysql")
	if err != nil {
		log.Fatal(err)
	}
}
