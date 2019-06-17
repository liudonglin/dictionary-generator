package handler

import (
	"code-server/store/base/db"
	"code-server/validator"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
)

type dbInfo struct {
	DBName string       `json:"name"`
	Tables []*tableInfo `json:"tables"`
}

type tableInfo struct {
	TableName string `json:"name"`
	Comments  string `json:"comments"`
}

func scanDB(scanner db.Scanner, dest *dbInfo) error {
	return scanner.Scan(
		&dest.DBName,
	)
}

func scanDBs(rows *sql.Rows) ([]*dbInfo, error) {
	defer rows.Close()

	dbs := []*dbInfo{}
	for rows.Next() {
		db := new(dbInfo)
		err := scanDB(rows, db)
		if err != nil {
			return nil, err
		}
		dbs = append(dbs, db)
	}
	return dbs, nil
}

func scanTable(scanner db.Scanner, dest *tableInfo) error {
	return scanner.Scan(
		&dest.TableName,
		&dest.Comments,
	)
}

func scanTables(rows *sql.Rows) ([]*tableInfo, error) {
	defer rows.Close()

	dbs := []*tableInfo{}
	for rows.Next() {
		db := new(tableInfo)
		err := scanTable(rows, db)
		if err != nil {
			return nil, err
		}
		dbs = append(dbs, db)
	}
	return dbs, nil
}

func loadConnInfo(c echo.Context) error {
	connInfo := &struct {
		DataBase string `json:"data_base" validate:"required"`
		HostPort string `json:"host_port" validate:"required"`
		User     string `json:"user" validate:"required"`
		Password string `json:"password" validate:"required"`
	}{}
	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, connInfo)

	//验证字段
	errs := validator.ValidateStruct(connInfo)
	if errs != nil {
		return &BusinessError{Message: errs.Error()}
	}

	if connInfo.DataBase != "mysql" {
		return &BusinessError{Message: "目前仅支持mysql数据库"}
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/information_schema", connInfo.User, connInfo.Password, connInfo.HostPort))
	if err != nil {
		return &BusinessError{Message: err.Error()}
	}

	err = db.Ping()
	if err != nil {
		return &BusinessError{Message: err.Error()}
	}

	conn := sqlx.NewDb(db, "mysql")
	defer conn.Close()

	dbs, _ := loadDBInfos(conn)

	return c.JSON(http.StatusOK, &StandardResult{
		Data: dbs,
	})
}

func loadDBInfos(conn *sqlx.DB) ([]*dbInfo, error) {
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
		tables, _ := loadTableInfos(conn, db.DBName)
		db.Tables = tables
	}

	return dbs, err
}

func loadTableInfos(conn *sqlx.DB, dbName string) ([]*tableInfo, error) {
	var tableSelectSQL = fmt.Sprintf(`SELECT 
	TABLE_NAME as tableName,
	TABLE_COMMENT as comments
	FROM information_schema.TABLES WHERE table_schema='%s'`, dbName)

	rows, err := conn.Query(tableSelectSQL)
	if err != nil {
		return nil, err
	}
	tables, err := scanTables(rows)
	return tables, err
}
