package handler

import (
	"code-server/core"
	"code-server/store"
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
	PID    int64        `json:"pid"`
	Tables []*tableInfo `json:"tables"`
}

type tableInfo struct {
	TableName string        `json:"name"`
	DBName    string        `json:"db_name"`
	Comment   string        `json:"comment"`
	Columns   []*columnInfo `json:"columns"`
}

type columnInfo struct {
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
}

type connInfo struct {
	DataBase string `json:"data_base" validate:"required"`
	HostPort string `json:"host_port" validate:"required"`
	User     string `json:"user" validate:"required"`
	Password string `json:"password" validate:"required"`
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
		&dest.DBName,
		&dest.Comment,
	)
}

func scanTables(rows *sql.Rows) ([]*tableInfo, error) {
	defer rows.Close()

	tabs := []*tableInfo{}
	for rows.Next() {
		tab := new(tableInfo)
		err := scanTable(rows, tab)
		if err != nil {
			return nil, err
		}
		tabs = append(tabs, tab)
	}
	return tabs, nil
}

func scanColumn(scanner db.Scanner, dest *columnInfo) error {
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
	)
}

func scanColumns(rows *sql.Rows) ([]*columnInfo, error) {
	defer rows.Close()

	cols := []*columnInfo{}
	for rows.Next() {
		col := new(columnInfo)
		err := scanColumn(rows, col)
		if err != nil {
			return nil, err
		}
		cols = append(cols, col)
	}
	return cols, nil
}

func loadConnInfo(c echo.Context) error {
	connInfo := &connInfo{}
	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, connInfo)

	//验证字段
	errs := validator.ValidateStruct(connInfo)
	if errs != nil {
		return &BusinessError{Message: errs.Error()}
	}

	conn, err := getConn(connInfo)
	if err != nil {
		return err
	}
	defer conn.Close()

	dbs, _ := loadDBInfos(conn)
	return c.JSON(http.StatusOK, &StandardResult{
		Data: dbs,
	})
}

func getConn(info *connInfo) (*sqlx.DB, error) {
	if info.DataBase != "mysql" {
		return nil, &BusinessError{Message: "目前仅支持mysql数据库"}
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/information_schema", info.User, info.Password, info.HostPort))
	if err != nil {
		return nil, &BusinessError{Message: err.Error()}
	}

	err = db.Ping()
	if err != nil {
		return nil, &BusinessError{Message: err.Error()}
	}

	conn := sqlx.NewDb(db, "mysql")
	return conn, nil
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
	TABLE_SCHEMA as dbName,
	TABLE_COMMENT as comments
	FROM information_schema.TABLES WHERE table_schema='%s'`, dbName)

	rows, err := conn.Query(tableSelectSQL)
	if err != nil {
		return nil, err
	}
	tables, err := scanTables(rows)
	for _, tab := range tables {
		columns, _ := loadColumnInfos(conn, tab.DBName, tab.TableName)
		tab.Columns = columns
	}
	return tables, err
}

func loadColumnInfos(conn *sqlx.DB, dbName, tableName string) ([]*columnInfo, error) {
	var columnSelectSQL = fmt.Sprintf(`SELECT 
	COLUMN_NAME as columnName,
	TABLE_NAME as tableName,
	TABLE_SCHEMA as dbName,
	COLUMN_COMMENT as comments,
	IS_NULLABLE as null_able,
	DATA_TYPE as data_type,
	COLUMN_TYPE as column_type,
	COLUMN_KEY as column_key,
	EXTRA as Extra
	FROM information_schema.COLUMNS WHERE table_schema='%s' And TABLE_NAME='%s'`, dbName, tableName)

	rows, err := conn.Query(columnSelectSQL)
	if err != nil {
		return nil, err
	}
	cols, err := scanColumns(rows)
	return cols, err
}

func saveConnInfo(c echo.Context) error {
	var dbInfos []*dbInfo
	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, &dbInfos)

	for _, dbPost := range dbInfos {
		//未选择table不保存
		if len(dbPost.Tables) == 0 {
			continue
		}
		dbStore := store.Stores().DataBaseStore
		//检查数据库是否存在，不存在则新增
		dbEntity, _ := dbStore.FindNameAndPID(dbPost.PID, dbPost.DBName)
		var dbID int64
		if dbEntity.ID > 0 {
			dbID = dbEntity.ID
		} else {
			dbInsert := &core.DataBase{
				Name: dbPost.DBName,
				PID:  dbPost.PID,
			}
			dbStore.Create(dbInsert)
			dbID = dbInsert.ID
		}

		for _, tablePost := range dbPost.Tables {
			tableStore := store.Stores().TableStore
			//检查数据表是否存在，不存在则新增
			tableEntity, _ := tableStore.FindNameAndDID(dbID, tablePost.TableName)
			var tableID int64
			if tableEntity.ID > 0 {
				tableID = tableEntity.ID
			} else {
				tableInsert := &core.Table{
					Name:        tablePost.TableName,
					Description: tablePost.Comment,
					DID:         dbID,
				}
				tableStore.Create(tableInsert)
				tableID = tableInsert.ID
			}

			for _, columnPost := range tablePost.Columns {
				columnStore := store.Stores().ColumnStore
				//检查数据列是否存在，不存在则新增
				columnEntity, _ := columnStore.FindNameAndTID(tableID, columnPost.ColumnName)
				if columnEntity.ID == 0 {
					columnInsert := &core.Column{
						Name:     columnPost.ColumnName,
						TID:      tableID,
						Title:    columnPost.Comment,
						DataType: columnPost.DataType,
					}
					columnStore.Create(columnInsert)
				}
			}
		}
	}

	return c.JSON(http.StatusOK, &StandardResult{
		Message: "保存成功!",
	})
}
