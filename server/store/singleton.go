package store

import (
	"sync"

	"../core"
	"../store/base/db"
	"./user"
)

// Singleton ...
type Singleton struct {
	UserStore core.UserStore
}

var singleton *Singleton
var once sync.Once

// GetInstance ...
func GetInstance(db *db.DB) *Singleton {
	once.Do(func() {
		userStore := user.New(db)
		singleton = &Singleton{userStore}
	})

	return singleton
}
