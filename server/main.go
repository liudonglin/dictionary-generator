package main

import (
	"fmt"
	"time"

	"dg-server/config"
	"dg-server/core"
	"dg-server/handler"
	"dg-server/store"
	"dg-server/store/base/db"
	"dg-server/tpl"

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
	// 初始化数据库操作类
	store.InitStores(_db)

	// 检查是否有初试账户,没有则创建
	if userCount, _ := store.Stores().UserStore.Count(); userCount == 0 {
		admin := &core.User{
			Login:    cfg.Admin.Username,
			Password: handler.EncryptionPassword(cfg.Admin.Password),
			Admin:    true,
			Active:   true,
			Created:  time.Now().Format("2006-01-02 15:04:05"),
		}
		store.Stores().UserStore.Create(admin)
	}

	port := 8080
	e := handler.New()

	tpl.TestGetTableScript(1)

	fmt.Printf("Access address http://localhost:%d/ui/#/dashboard \n", port)
	e.Start(fmt.Sprintf(":%d", port))
}
