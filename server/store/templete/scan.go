package templete

import (
	"database/sql"
	"dg-server/core"
	"dg-server/store/base/db"
)

// helper function converts the User structure to a set
// of named query parameters.
func toParams(u *core.Templete) map[string]interface{} {
	return map[string]interface{}{
		"templete_id":        u.ID,
		"templete_name":      u.Name,
		"templete_content":   u.Content,
		"templete_language":  u.Language,
		"templete_data_base": u.DataBase,
		"templete_orm":       u.Orm,
		"templete_type":      u.Type,
		"templete_created":   u.Created,
		"templete_updated":   u.Updated,
	}
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRow(scanner db.Scanner, dest *core.Templete) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.Content,
		&dest.Language,
		&dest.DataBase,
		&dest.Orm,
		&dest.Type,
		&dest.Created,
		&dest.Updated,
	)
}

func scanSingle(scanner db.Scanner, dest interface{}) error {
	return scanner.Scan(
		dest,
	)
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRows(rows *sql.Rows) ([]*core.Templete, error) {
	defer rows.Close()

	tpls := []*core.Templete{}
	for rows.Next() {
		tpl := new(core.Templete)
		err := scanRow(rows, tpl)
		if err != nil {
			return nil, err
		}
		tpls = append(tpls, tpl)
	}
	return tpls, nil
}
