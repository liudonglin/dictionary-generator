package store

import (
	"crypto/md5"
	"fmt"
	"sync"
	"time"

	"dg-server/config"
	"dg-server/core"
	"dg-server/store/base/db"
	"dg-server/store/column"
	"dg-server/store/connection"
	"dg-server/store/dbase"
	"dg-server/store/project"
	"dg-server/store/table"
	"dg-server/store/templete"
	"dg-server/store/user"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

// Singleton ...
type Singleton struct {
	UserStore       core.UserStore
	ProjectStore    core.ProjectStore
	DataBaseStore   core.DataBaseStore
	TableStore      core.TableStore
	ColumnStore     core.ColumnStore
	ConnectionStore core.ConnectionStore
	TempleteStore   core.TempleteStore
}

var singleton *Singleton
var once sync.Once

// Stores ...
func Stores() *Singleton {
	return singleton
}

func init() {
	cfg, err := config.Environ()
	if err != nil {
		logger := logrus.WithError(err)
		logger.Fatalln("main: invalid configuration")
	}

	once.Do(func() {
		_db, err := db.Connect(cfg.Database.Driver, cfg.Database.Datasource)
		if err != nil {
			panic(err)
		}
		singleton = &Singleton{
			UserStore:       user.New(_db),
			ProjectStore:    project.New(_db),
			DataBaseStore:   dbase.New(_db),
			TableStore:      table.New(_db),
			ColumnStore:     column.New(_db),
			ConnectionStore: connection.New(_db),
			TempleteStore:   templete.New(_db),
		}
	})

	// 检查是否有初试账户,没有则创建
	if userCount, _ := Stores().UserStore.Count(); userCount == 0 {
		admin := &core.User{
			Login:    cfg.Admin.Username,
			Password: EncryptionPassword(cfg.Admin.Password),
			Admin:    true,
			Active:   true,
			Created:  time.Now().Format("2006-01-02 15:04:05"),
		}
		Stores().UserStore.Create(admin)
	}
}

// EncryptionPassword md5加密用户密码
func EncryptionPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", hash)
}
