package imports

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func loadMysqlDBInfos(conn *sqlx.DB) ([]*DBInfo, error) {
	var dbSelectSQL = `select schema_name from information_schema.schemata 
	where schema_name !='information_schema'
	And schema_name !='performance_schema'
	And schema_name !='mysql'
	And schema_name !='sys'
	`
	rows, err := conn.Query(dbSelectSQL)
	if err != nil {
		return nil, err
	}
	dbs, err := scanDBs(rows)

	for _, db := range dbs {
		tables, _ := loadMysqlTableInfos(conn, db.DBName)
		db.Tables = tables
	}

	return dbs, err
}

func loadMysqlTableInfos(conn *sqlx.DB, dbName string) ([]*TableInfo, error) {
	var tableSelectSQL = fmt.Sprintf(`SELECT 
	TABLE_NAME as tableName,
	TABLE_SCHEMA as dbName,
	TABLE_COMMENT as comments
	FROM information_schema.TABLES WHERE table_schema='%s'`, dbName)

	rows, err := conn.Query(tableSelectSQL)
	if err != nil {
		return nil, err
	}
	tables, err := scanTables(rows)
	for _, tab := range tables {
		columns, _ := loadMysqlColumnInfos(conn, tab.DBName, tab.TableName)
		tab.Columns = columns
	}
	return tables, err
}

func loadMysqlColumnInfos(conn *sqlx.DB, dbName, tableName string) ([]*ColumnInfo, error) {
	var columnSelectSQL = fmt.Sprintf(`SELECT 
	COLUMN_NAME as columnName,
	TABLE_NAME as tableName,
	TABLE_SCHEMA as dbName,
	COLUMN_COMMENT as comments,
	IS_NULLABLE as null_able,
	DATA_TYPE as data_type,
	COLUMN_TYPE as column_type,
	COLUMN_KEY as column_key,
	EXTRA as Extra,
	CHARACTER_MAXIMUM_LENGTH as char_length,
	NUMERIC_PRECISION as num_precision,
	NUMERIC_SCALE as num_scale
	FROM information_schema.COLUMNS WHERE table_schema='%s' And TABLE_NAME='%s'`, dbName, tableName)

	rows, err := conn.Query(columnSelectSQL)
	if err != nil {
		return nil, err
	}
	cols, err := scanColumns(rows)
	return cols, err
}
