package imports

import (
	"code-server/store/base/db"
	"database/sql"
)

// DBInfo ...
type DBInfo struct {
	DBName string       `json:"name"`
	PID    int64        `json:"pid"`
	Tables []*TableInfo `json:"tables"`
}

// TableInfo ...
type TableInfo struct {
	TableName string        `json:"name"`
	DBName    string        `json:"db_name"`
	Comment   string        `json:"comment"`
	Columns   []*ColumnInfo `json:"columns"`
}

// ColumnInfo ...
type ColumnInfo struct {
	ColumnName string `json:"name"`
	TableName  string `json:"table_name"`
	DBName     string `json:"db_name"`
	Comment    string `json:"comment"`
	Null       string `json:"null"`
	DataType   string `json:"data_type"`
	ColumnType string `json:"column_type"`
	//索引类型。可包含的值有PRI，代表主键，UNI，代表唯一键，MUL，可重复
	ColumnKey string `json:"column_key"`
	Extra     string `json:"extra"`
	//字符最大长度
	CharacterMaximumLength sql.NullString `json:"char_length"`
	//整数长度
	NumericPrecision sql.NullString `json:"num_precision"`
	//小数长度
	NumericScale sql.NullString `json:"num_scale"`
}

func scanDB(scanner db.Scanner, dest *DBInfo) error {
	return scanner.Scan(
		&dest.DBName,
	)
}

func scanDBs(rows *sql.Rows) ([]*DBInfo, error) {
	defer rows.Close()

	dbs := []*DBInfo{}
	for rows.Next() {
		db := new(DBInfo)
		err := scanDB(rows, db)
		if err != nil {
			return nil, err
		}
		dbs = append(dbs, db)
	}
	return dbs, nil
}

func scanTable(scanner db.Scanner, dest *TableInfo) error {
	return scanner.Scan(
		&dest.TableName,
		&dest.DBName,
		&dest.Comment,
	)
}

func scanTables(rows *sql.Rows) ([]*TableInfo, error) {
	defer rows.Close()

	tabs := []*TableInfo{}
	for rows.Next() {
		tab := new(TableInfo)
		err := scanTable(rows, tab)
		if err != nil {
			return nil, err
		}
		tabs = append(tabs, tab)
	}
	return tabs, nil
}

func scanColumn(scanner db.Scanner, dest *ColumnInfo) error {
	return scanner.Scan(
		&dest.ColumnName,
		&dest.TableName,
		&dest.DBName,
		&dest.Comment,
		&dest.Null,
		&dest.DataType,
		&dest.ColumnType,
		&dest.ColumnKey,
		&dest.Extra,
		&dest.CharacterMaximumLength,
		&dest.NumericPrecision,
		&dest.NumericScale,
	)
}

func scanColumns(rows *sql.Rows) ([]*ColumnInfo, error) {
	defer rows.Close()

	cols := []*ColumnInfo{}
	for rows.Next() {
		col := new(ColumnInfo)
		err := scanColumn(rows, col)
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
	}
	return cols, nil
}
