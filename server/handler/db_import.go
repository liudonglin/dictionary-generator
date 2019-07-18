package handler

import (
	"dg-server/core"
	"dg-server/imports"
	"dg-server/store"
	"dg-server/validator"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func loadDBsByConnInfo(c echo.Context) error {
	connInfo := &core.Connection{}
	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, connInfo)

	//验证字段
	errs := validator.ValidateStruct(connInfo)
	if errs != nil {
		return &BusinessError{Message: errs.Error()}
	}

	dbs, err := imports.LoadDBsByConnInfo(connInfo)
	if err != nil {
		return &BusinessError{Message: err.Error()}
	}

	return c.JSON(http.StatusOK, &StandardResult{
		Data: dbs,
	})
}

func saveConnInfo(c echo.Context) error {
	var dbInfos []*imports.DBInfo
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
			err := dbStore.Create(dbInsert)
			if err != nil {
				return err
			}
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
					PID:         dbPost.PID,
				}
				err := tableStore.Create(tableInsert)
				if err != nil {
					return err
				}
				tableID = tableInsert.ID
			}

			for _, columnPost := range tablePost.Columns {
				columnStore := store.Stores().ColumnStore
				//检查数据列是否存在，不存在则新增
				columnEntity, _ := columnStore.FindNameAndTID(tableID, columnPost.ColumnName)
				if columnEntity.ID == 0 {

					pk, index, unique, ai, null, length := false, false, false, false, false, ""
					if columnPost.ColumnKey == "PRI" {
						//主键索引
						pk = true
						index = true
						unique = true
					}
					if columnPost.ColumnKey == "UNI" {
						//唯一索引
						index = true
						unique = true
					}
					if columnPost.ColumnKey == "MUL" {
						//一般索引
						index = true
					}

					if columnPost.Extra == "auto_increment" {
						ai = true
					}
					if columnPost.Null == "YES" {
						null = true
					}

					if columnPost.DataType == "int" ||
						columnPost.DataType == "tinyint" ||
						columnPost.DataType == "smallint" ||
						columnPost.DataType == "integer" ||
						columnPost.DataType == "bigint" ||
						columnPost.DataType == "bit" {
						if columnPost.NumericPrecision.Valid {
							length = columnPost.NumericPrecision.String
						}
					}
					if columnPost.DataType == "float" ||
						columnPost.DataType == "decimal" {

						if columnPost.NumericPrecision.Valid {
							length = columnPost.NumericPrecision.String + ","
						}
						if columnPost.NumericScale.Valid {
							length += columnPost.NumericScale.String
						} else {
							length += "0"
						}
					}
					if strings.Contains(columnPost.DataType, "char") {
						if columnPost.CharacterMaximumLength.Valid {
							length = columnPost.CharacterMaximumLength.String
						}
					}

					columnInsert := &core.Column{
						Name:       columnPost.ColumnName,
						PID:        dbPost.PID,
						DID:        dbID,
						TID:        tableID,
						Title:      columnPost.Comment,
						DataType:   columnPost.DataType,
						ColumnType: columnPost.ColumnType,
						PK:         pk,
						Unique:     unique,
						AI:         ai,
						Null:       null,
						Index:      index,
						Length:     length,
					}
					err := columnStore.Create(columnInsert)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return c.JSON(http.StatusOK, &StandardResult{
		Message: "保存成功!",
	})
}
