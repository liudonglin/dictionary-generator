package database

import (
	"code-server/core"
	"code-server/store/base/db"
	"database/sql"
)

// helper function converts the User structure to a set
// of named query parameters.
func toParams(u *core.DataBase) map[string]interface{} {
	return map[string]interface{}{
		"database_id":          u.ID,
		"database_name":        u.Name,
		"database_pid":         u.PID,
		"database_description": u.Description,
		"database_created":     u.Created,
		"database_updated":     u.Updated,
	}
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRow(scanner db.Scanner, dest *core.DataBase) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.PID,
		&dest.Description,
		&dest.Created,
		&dest.Updated,
	)
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRows(rows *sql.Rows) ([]*core.DataBase, error) {
	defer rows.Close()

	dbs := []*core.DataBase{}
	for rows.Next() {
		db := new(core.DataBase)
		err := scanRow(rows, db)
		if err != nil {
			return nil, err
		}
		dbs = append(dbs, db)
	}
	return dbs, nil
}
