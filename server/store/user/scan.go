package user

import (
	"database/sql"

	"code-server/core"
	"code-server/store/base/db"
)

// helper function converts the User structure to a set
// of named query parameters.
func toParams(u *core.User) map[string]interface{} {
	return map[string]interface{}{
		"user_id":         u.ID,
		"user_login":      u.Login,
		"user_password":   u.Password,
		"user_email":      u.Email,
		"user_admin":      u.Admin,
		"user_active":     u.Active,
		"user_avatar":     u.Avatar,
		"user_created":    u.Created,
		"user_updated":    u.Updated,
		"user_last_login": u.LastLogin,
	}
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRow(scanner db.Scanner, dest *core.User) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Login,
		&dest.Password,
		&dest.Email,
		&dest.Admin,
		&dest.Active,
		&dest.Avatar,
		&dest.Created,
		&dest.Updated,
		&dest.LastLogin,
	)
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRows(rows *sql.Rows) ([]*core.User, error) {
	defer rows.Close()

	users := []*core.User{}
	for rows.Next() {
		user := new(core.User)
		err := scanRow(rows, user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
