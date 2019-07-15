package store

import (
	"sync"

	"dg-server/core"
	"dg-server/store/base/db"
	"dg-server/store/column"
	"dg-server/store/connection"
	"dg-server/store/dbase"
	"dg-server/store/project"
	"dg-server/store/table"
	"dg-server/store/templete"
	"dg-server/store/user"
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
	if singleton == nil {
		panic("has not init stores !")
	}
	return singleton
}

// InitStores ...
func InitStores(db *db.DB) {
	once.Do(func() {
		singleton = &Singleton{
			UserStore:       user.New(db),
			ProjectStore:    project.New(db),
			DataBaseStore:   dbase.New(db),
			TableStore:      table.New(db),
			ColumnStore:     column.New(db),
			ConnectionStore: connection.New(db),
			TempleteStore:   templete.New(db),
		}
	})
}
