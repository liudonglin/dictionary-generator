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

					pk, index, ai, null, length := false, false, false, false, ""
					if columnPost.ColumnKey == "PRI" {
						pk = true
					}
					if columnPost.ColumnKey == "UNI" {
						index = true
					}
					// MUL 组合索引不支持
					if columnPost.Extra == "auto_increment" {
						ai = true
					}
					if columnPost.Null == "YES" {
						null = true
					}

					if strings.Contains(columnPost.DataType, "int") ||
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
					if strings.Contains(columnPost.DataType, "char") ||
						columnPost.DataType == "bit" {
						if columnPost.CharacterMaximumLength.Valid {
							length = columnPost.CharacterMaximumLength.String
						}
					}

					columnInsert := &core.Column{
						Name:     columnPost.ColumnName,
						TID:      tableID,
						Title:    columnPost.Comment,
						DataType: columnPost.DataType,
						PK:       pk,
						AI:       ai,
						Null:     null,
						Index:    index,
						Length:   length,
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
