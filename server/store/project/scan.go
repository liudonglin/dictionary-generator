package project

import (
	"database/sql"

	"dg-server/core"
	"dg-server/store/base/db"
)

// helper function converts the User structure to a set
// of named query parameters.
func toParams(u *core.Project) map[string]interface{} {
	return map[string]interface{}{
		"project_id":          u.ID,
		"project_name":        u.Name,
		"project_language":    u.Language,
		"project_data_base":   u.DataBase,
		"project_orm":         u.Orm,
		"project_description": u.Description,
		"project_name_space":  u.NameSpace,
		"project_created":     u.Created,
		"project_updated":     u.Updated,
	}
}

// helper function scans the sql.Row and copies the column
// values to the destination object.
func scanRow(scanner db.Scanner, dest *core.Project) error {
	return scanner.Scan(
		&dest.ID,
		&dest.Name,
		&dest.Language,
		&dest.DataBase,
		&dest.Orm,
		&dest.Description,
		&dest.NameSpace,
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
func scanRows(rows *sql.Rows) ([]*core.Project, error) {
	defer rows.Close()

	projects := []*core.Project{}
	for rows.Next() {
		project := new(core.Project)
		err := scanRow(rows, project)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, nil
}
