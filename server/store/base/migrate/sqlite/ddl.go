package sqlite

import (
	"database/sql"
)

var migrations = []struct {
	name string
	stmt string
}{
	{
		name: "create-table-users",
		stmt: createTableUsers,
	},
	{
		name: "create-table-projects",
		stmt: createTableProjects,
	},
	{
		name: "create-table-database",
		stmt: createTableDatabBses,
	},
	{
		name: "create-table-tables",
		stmt: createTableTables,
	},
	{
		name: "create-table-columns",
		stmt: createTableColumns,
	},
}

// Migrate performs the database migration. If the migration fails
// and error is returned.
func Migrate(db *sql.DB) error {
	if err := createTable(db); err != nil {
		return err
	}
	completed, err := selectCompleted(db)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	for _, migration := range migrations {
		if _, ok := completed[migration.name]; ok {

			continue
		}

		if _, err := db.Exec(migration.stmt); err != nil {
			return err
		}
		if err := insertMigration(db, migration.name); err != nil {
			return err
		}

	}
	return nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(migrationTableCreate)
	return err
}

func insertMigration(db *sql.DB, name string) error {
	_, err := db.Exec(migrationInsert, name)
	return err
}

func selectCompleted(db *sql.DB) (map[string]struct{}, error) {
	migrations := map[string]struct{}{}
	rows, err := db.Query(migrationSelect)
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

//
// migration table ddl and sql
//

var migrationTableCreate = `
CREATE TABLE IF NOT EXISTS migrations (
 name VARCHAR(255)
,UNIQUE(name)
)
`

var migrationInsert = `
INSERT INTO migrations (name) VALUES (?)
`

var migrationSelect = `
SELECT name FROM migrations
`

//
// 001_create_table_user.sql
//

var createTableUsers = `
CREATE TABLE IF NOT EXISTS users (
 user_id            INTEGER PRIMARY KEY AUTOINCREMENT
,user_login         TEXT
,user_password      TEXT
,user_email         TEXT
,user_admin         BOOLEAN
,user_active        BOOLEAN
,user_avatar        TEXT
,user_created       TEXT
,user_updated       TEXT
,user_last_login    TEXT
,UNIQUE(user_login)
);
`

var createTableProjects = `
CREATE TABLE IF NOT EXISTS projects (
	project_id            	INTEGER PRIMARY KEY AUTOINCREMENT
   ,project_name         	TEXT
   ,project_language      	TEXT
   ,project_data_base       TEXT
   ,project_orm         	TEXT
   ,project_description     TEXT
   ,project_created       	TEXT
   ,project_updated       	TEXT
   ,UNIQUE(project_name)
   );
`

var createTableDatabBses = `
CREATE TABLE IF NOT EXISTS database (
	database_id            	INTEGER PRIMARY KEY AUTOINCREMENT
   ,database_name         	TEXT
   ,database_pid      		INTEGER
   ,database_description    TEXT
   ,database_created       	TEXT
   ,database_updated       	TEXT
   ,UNIQUE(database_pid, database_name)
   );
`

var createTableTables = `
CREATE TABLE IF NOT EXISTS tables (
	table_id            	INTEGER PRIMARY KEY AUTOINCREMENT
   ,table_name         		TEXT
   ,table_did      			INTEGER
   ,table_description    	TEXT
   ,table_created       	TEXT
   ,table_updated       	TEXT
   ,UNIQUE(table_did, table_name)
   );
`

var createTableColumns = `
CREATE TABLE IF NOT EXISTS columns (
	column_id            	INTEGER PRIMARY KEY AUTOINCREMENT
   ,column_name         	TEXT
   ,column_tid      		INTEGER
   ,column_title			TEXT
   ,column_data_type		TEXT
   ,column_pk				INTEGER
   ,column_ai				INTEGER
   ,column_null				INTEGER
   ,column_length			TEXT
   ,column_index			INTEGER
   ,column_enum    			TEXT
   ,column_description    	TEXT
   ,column_created       	TEXT
   ,column_updated       	TEXT
   ,UNIQUE(column_tid, column_name)
);
`
