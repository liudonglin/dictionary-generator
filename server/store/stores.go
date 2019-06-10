package store

import (
	"sync"

	"../core"
	"../store/base/db"
	"./project"
	"./user"
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
