package store

import (
	"sync"

	"code-server/core"
	"code-server/store/base/db"
	"code-server/store/project"
	"code-server/store/user"
)

// Singleton ...
type Singleton struct {
	UserStore    core.UserStore
	ProjectStore core.ProjectStore
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
			UserStore:    user.New(db),
			ProjectStore: project.New(db),
		}
	})
}
