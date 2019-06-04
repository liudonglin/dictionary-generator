package user

import (
	"context"

	"../../core"
	"../base/db"
)

// New returns a new UserStore.
func New(db *db.DB) core.UserStore {
	return &userStore{db}
}

type userStore struct {
	db *db.DB
}

// Create persists a new user to the datastore.
func (s *userStore) Create(ctx context.Context, user *core.User) error {
	return s.db.Lock(func(execer db.Execer, binder db.Binder) error {
		params := toParams(user)
		stmt, args, err := binder.BindNamed(stmtInsert, params)
		if err != nil {
			return err
		}
		res, err := execer.Exec(stmt, args...)
		if err != nil {
			return err
		}
		user.ID, err = res.LastInsertId()
		return err
	})
}

const stmtInsert = `
INSERT INTO users (
 user_login
,user_email
,user_admin
,user_active
,user_avatar
,user_created
,user_updated
,user_last_login
) VALUES (
 :user_login
,:user_email
,:user_admin
,:user_active
,:user_avatar
,:user_created
,:user_updated
,:user_last_login
)
`
