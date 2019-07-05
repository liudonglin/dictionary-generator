package column

import (
	"database/sql"

	"dg-server/core"
	"dg-server/store/base/db"
)

// helper function converts the User structure to a set
// of named query parameters.
func toParams(u *core.Column) map[string]interface{} {
	return map[string]interface{}{
		"column_id":          u.ID,
		"column_name":        u.Name,
		"column_pid":         u.PID,
		"column_did":         u.DID,
		"column_tid":         u.TID,
		"column_title":       u.Title,
		"column_data_type":   u.DataType,
		"column_pk":          u.PK,
		"column_ai":          u.AI,
		"column_null":        u.Null,
		"column_length":      u.Length,
		"column_index":       u.Index,
		"column_enum":        u.Enum,
		"column_description": u.Description,
		"column_created":     u.Created,
		"column_updated":     u.Updated,
	}
}

func scanSingle(scanner db.Scanner, dest interface{}) error {
	return scanner.Scan(
		dest,
	)
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRow(scanner db.Scanner, dest *core.Column) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.PID,
		&dest.DID,
		&dest.TID,
		&dest.Title,
		&dest.DataType,
		&dest.PK,
		&dest.AI,
		&dest.Null,
		&dest.Length,
		&dest.Index,
		&dest.Enum,
		&dest.Description,
		&dest.Created,
		&dest.Updated,
	)
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRows(rows *sql.Rows) ([]*core.Column, error) {
	defer rows.Close()

	columns := []*core.Column{}
	for rows.Next() {
		column := new(core.Column)
		err := scanRow(rows, column)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}
	return columns, nil
}
