package main

import (
	"fmt"

	"./config"
	"./handler"
	"./store/base/db"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func main() {

	cfg, err := config.Environ()
	if err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("main: invalid configuration")
	}

	_db, err := db.Connect(cfg.Database.Driver, cfg.Database.Datasource)
	if err != nil {
		panic(err)
	}

	port := 8080
	e := handler.New(_db)

	fmt.Printf("Access address http://localhost:%d/ui/#/dashboard \n", port)
	e.Start(fmt.Sprintf(":%d", port))
}
