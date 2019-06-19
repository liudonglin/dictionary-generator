package imports

import (
	"code-server/core"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

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
