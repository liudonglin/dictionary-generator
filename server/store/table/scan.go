package table

import (
	"database/sql"

	"dg-server/core"
	"dg-server/store/base/db"
)

// helper function converts the User structure to a set
// of named query parameters.
func toParams(u *core.Table) map[string]interface{} {
	return map[string]interface{}{
		"table_id":          u.ID,
		"table_name":        u.Name,
		"table_pid":         u.PID,
		"table_did":         u.DID,
		"table_title":       u.Title,
		"table_description": u.Description,
		"table_created":     u.Created,
		"table_updated":     u.Updated,
	}
}

func scanSingle(scanner db.Scanner, dest interface{}) error {
	return scanner.Scan(
		dest,
	)
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRow(scanner db.Scanner, dest *core.Table) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.PID,
		&dest.DID,
		&dest.Title,
		&dest.Description,
		&dest.Created,
		&dest.Updated,
	)
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRows(rows *sql.Rows) ([]*core.Table, error) {
	defer rows.Close()

	tables := []*core.Table{}
	for rows.Next() {
		table := new(core.Table)
		err := scanRow(rows, table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}
