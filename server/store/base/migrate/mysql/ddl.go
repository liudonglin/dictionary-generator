package mysql

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
	{
		name: "create-table-connections",
		stmt: createTableConnections,
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
 user_id            INTEGER PRIMARY KEY AUTO_INCREMENT
,user_login         VARCHAR(50)
,user_password      VARCHAR(50)
,user_email         VARCHAR(200)
,user_admin         BOOLEAN
,user_active        BOOLEAN
,user_avatar        VARCHAR(500)
,user_created       VARCHAR(20)
,user_updated       VARCHAR(20)
,user_last_login    VARCHAR(20)
,UNIQUE(user_login)
);
`

var createTableProjects = `
CREATE TABLE IF NOT EXISTS projects (
	project_id            	INTEGER PRIMARY KEY AUTO_INCREMENT
   ,project_name         	VARCHAR(20)
   ,project_language      	VARCHAR(10)
   ,project_data_base       VARCHAR(10)
   ,project_orm         	VARCHAR(10)
   ,project_description     VARCHAR(200)
   ,project_created       	VARCHAR(20)
   ,project_updated       	VARCHAR(20)
   ,UNIQUE(project_name)
   );
`

var createTableDatabBses = `
CREATE TABLE IF NOT EXISTS database (
	database_id            	INTEGER PRIMARY KEY AUTO_INCREMENT
   ,database_name         	VARCHAR(20)
   ,database_pid      		INTEGER
   ,database_description    VARCHAR(200)
   ,database_created       	VARCHAR(20)
   ,database_updated       	VARCHAR(20)
   ,UNIQUE(database_pid, database_name)
   );
`

var createTableTables = `
CREATE TABLE IF NOT EXISTS tables (
	table_id            	INTEGER PRIMARY KEY AUTO_INCREMENT
   ,table_name         		VARCHAR(20)
   ,table_did      			INTEGER
   ,table_description    	VARCHAR(200)
   ,table_created       	VARCHAR(20)
   ,table_updated       	VARCHAR(20)
   ,UNIQUE(table_did, table_name)
   );
`

var createTableColumns = `
CREATE TABLE IF NOT EXISTS columns (
	column_id            	INTEGER PRIMARY KEY AUTO_INCREMENT
   ,column_name         	VARCHAR(20)
   ,column_tid      		INTEGER
   ,column_title			VARCHAR(20)
   ,column_data_type		VARCHAR(10)
   ,column_pk				INTEGER
   ,column_ai				INTEGER
   ,column_null				INTEGER
   ,column_length			VARCHAR(10)
   ,column_index			INTEGER
   ,column_enum    			VARCHAR(100)
   ,column_description    	VARCHAR(200)
   ,column_created       	VARCHAR(20)
   ,column_updated       	VARCHAR(20)
   ,UNIQUE(column_tid, column_name)
);
`

var createTableConnections = `
CREATE TABLE IF NOT EXISTS connections (
	connection_id            	INTEGER PRIMARY KEY AUTO_INCREMENT
   ,connection_name         	VARCHAR(20)
   ,connection_pid      		INTEGER
   ,connection_data_base		VARCHAR(20)
   ,connection_host				VARCHAR(20)
   ,connection_port				INTEGER
   ,connection_user				INTEGER
   ,connection_password			INTEGER
   ,connection_description    	VARCHAR(200)
   ,connection_created       	VARCHAR(20)
   ,connection_updated       	VARCHAR(20)
   ,UNIQUE(connection_pid, connection_name)
);
`
