package imports

import (
	"dg-server/core"
	"dg-server/store/base/db"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

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

func getConn(info *core.Connection) (*sqlx.DB, error) {

	if info.DataBase != "mysql" {
		return nil, errors.New("目前仅支持mysql")
	}

	connString := ""
	if info.DataBase == "mssql" {
		connString = fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=%s;encrypt=disable", info.Host, "database", info.User, info.Password, info.Port)
	} else if info.DataBase == "mysql" {
		connString = fmt.Sprintf("%s:%s@tcp(%s:%s)/information_schema", info.User, info.Password, info.Host, info.Port)
	} else {
		return nil, errors.New("不支持的数据库类型:" + info.DataBase)
	}

	db, err := sql.Open(info.DataBase, connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	conn := sqlx.NewDb(db, info.DataBase)
	return conn, nil
}

// LoadDBsByConnInfo ...
func LoadDBsByConnInfo(info *core.Connection) ([]*DBInfo, error) {
	conn, err := getConn(info)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if info.DataBase == "mysql" {
		return loadMysqlDBInfos(conn)
	}

	return nil, errors.New("不支持的数据库类型:" + info.DataBase)
}
