package handler

import (
	"dg-server/core"
	"dg-server/store"
	"dg-server/validator"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func saveTable(c echo.Context) error {
	table := &core.Table{}
	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, table)

	//验证字段
	errs := validator.ValidateStruct(table)
	if errs != nil {
		return &BusinessError{Message: errs.Error()}
	}

	tableStore := store.Stores().TableStore

	//检查表名称是否重复
	dbTable, err := tableStore.FindNameAndDID(table.DID, table.Name)
	if err != nil {
		return err
	}
	if dbTable.ID != 0 && dbTable.ID != table.ID {
		return &BusinessError{Message: fmt.Sprintf("名称为%s的表已存在!", table.Name)}
	}

	if table.ID == 0 {
		err := tableStore.Create(table)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, &StandardResult{
			Message: "新增成功!",
			Data:    table.ID,
		})
	}

	err = tableStore.Update(table)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &StandardResult{
		Message: "修改成功!",
		Data:    table.ID,
	})
}

func deleteTable(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)
	id, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		return err
	}
	tableStore := store.Stores().TableStore
	err = tableStore.Delete(id)
	if err != nil {
		return err
	}

	columnStore := store.Stores().ColumnStore
	err = columnStore.DeleteByTID(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &StandardResult{
		Message: "删除成功",
	})
}
