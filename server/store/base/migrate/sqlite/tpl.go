package sqlite

import (
	"database/sql"
	"dg-server/core"
	"dg-server/store/base/migrate"
	"time"
)

// MigrateTemplete 初始化Templete
func initTemplete(db *sql.DB) error {
	completed, _ := selectCompletedTpl(db)
	for _, mtpl := range migrate.MigrationTempletes {

		if _, ok := completed[mtpl.Name]; ok {
			continue
		}

		if err := insertTemplete(db, &mtpl); err != nil {
			return err
		}
	}
	return nil
}

func insertTemplete(db *sql.DB, tpl *core.Templete) error {
	created := time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Exec(migrationTempleteInsert,
		tpl.Name,
		tpl.Content,
		tpl.Language,
		tpl.DataBase,
		tpl.Orm,
		tpl.Type,
		created,
		created)
	return err
}

func selectCompletedTpl(db *sql.DB) (map[string]struct{}, error) {
	migrations := map[string]struct{}{}
	rows, err := db.Query(migrationSelectTpl)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		migrations[name] = struct{}{}
	}
	return migrations, nil
}

const migrationTempleteInsert = `
INSERT INTO templetes (
templete_name
,templete_content
,templete_language
,templete_data_base
,templete_orm
,templete_type
,templete_created
,templete_updated
) VALUES (
?
,?
,?
,?
,?
,?
,?
,?
)
`

var migrationSelectTpl = `
SELECT templete_name FROM templetes
`
