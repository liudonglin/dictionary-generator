package connection

import (
	"code-server/core"
	"code-server/store/base/db"
	"database/sql"
)

// helper function converts the User structure to a set
// of named query parameters.
func toParams(u *core.Connection) map[string]interface{} {
	return map[string]interface{}{
		"connection_id":          u.ID,
		"connection_name":        u.Name,
		"connection_pid":         u.PID,
		"connection_data_base":   u.DataBase,
		"connection_host":        u.Host,
		"connection_port":        u.Port,
		"connection_user":        u.User,
		"connection_password":    u.Password,
		"connection_description": u.Description,
		"connection_created":     u.Created,
		"connection_updated":     u.Updated,
	}
}

func scanSingle(scanner db.Scanner, dest interface{}) error {
	return scanner.Scan(
		dest,
	)
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRow(scanner db.Scanner, dest *core.Connection) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.PID,
		&dest.DataBase,
		&dest.Host,
		&dest.Port,
		&dest.User,
		&dest.Password,
		&dest.Description,
		&dest.Created,
		&dest.Updated,
	)
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRows(rows *sql.Rows) ([]*core.Connection, error) {
	defer rows.Close()

	connections := []*core.Connection{}
	for rows.Next() {
		connection := new(core.Connection)
		err := scanRow(rows, connection)
		if err != nil {
			return nil, err
		}
		connections = append(connections, connection)
	}
	return connections, nil
}
