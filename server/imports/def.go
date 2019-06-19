package imports

import (
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
