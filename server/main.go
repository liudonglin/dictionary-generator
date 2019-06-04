package main

import (
	"fmt"

	"./config"
	"./core"
	"./handler"
	"./store"
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

	// 检查是否有初试账户
	instance := store.GetInstance(_db)
	if userCount, _ := instance.UserStore.Count(); userCount == 0 {
		admin := &core.User{
			Login:    cfg.Admin.Username,
			Password: cfg.Admin.Password,
			Admin:    true,
			Active:   true,
		}
		instance.UserStore.Create(admin)
	}

	port := 8080
	e := handler.New()

	fmt.Printf("Access address http://localhost:%d/ui/#/dashboard \n", port)
	e.Start(fmt.Sprintf(":%d", port))
}
