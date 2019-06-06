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

// GetInstance ...
func GetInstance(db *db.DB) *Singleton {
	once.Do(func() {
		singleton = &Singleton{
			UserStore:    user.New(db),
			ProjectStore: project.New(db),
		}
	})

	return singleton
}
