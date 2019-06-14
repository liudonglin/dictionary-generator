package store

import (
	"sync"

	"code-server/core"
	"code-server/store/base/db"
	"code-server/store/column"
	"code-server/store/dbase"
	"code-server/store/project"
	"code-server/store/table"
	"code-server/store/user"
)

// Singleton ...
type Singleton struct {
	UserStore     core.UserStore
	ProjectStore  core.ProjectStore
	DataBaseStore core.DataBaseStore
	TableStore    core.TableStore
	ColumnStore   core.ColumnStore
}

var singleton *Singleton
var once sync.Once

// Stores ...
func Stores() *Singleton {
	return singleton
}

// InitStores ...
func InitStores(db *db.DB) {
	once.Do(func() {
		singleton = &Singleton{
			UserStore:     user.New(db),
			ProjectStore:  project.New(db),
			DataBaseStore: dbase.New(db),
			TableStore:    table.New(db),
			ColumnStore:   column.New(db),
		}
	})
}
