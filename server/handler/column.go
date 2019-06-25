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

func saveColumn(c echo.Context) error {
	postEntity := &core.Column{}
	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, postEntity)

	//验证字段
	errs := validator.ValidateStruct(postEntity)
	if errs != nil {
		return &BusinessError{Message: errs.Error()}
	}

	cStore := store.Stores().ColumnStore

	//检查名称是否重复
	dbEntity, err := cStore.FindNameAndTID(postEntity.TID, postEntity.Name)
	if err != nil {
		return err
	}
	if dbEntity.ID != 0 && dbEntity.ID != postEntity.ID {
		return &BusinessError{Message: fmt.Sprintf("名称为%s的列在表中已存在!", postEntity.Name)}
	}

	//检查主键是否重复
	if postEntity.PK {
		pkEntity, err := cStore.FindPK(postEntity.TID)
		if err != nil {
			return err
		}
		if pkEntity.ID != 0 && pkEntity.ID != postEntity.ID {
			return &BusinessError{Message: fmt.Sprintf("表中已存在名为%s的主键!", pkEntity.Name)}
		}
	}

	//add
	if postEntity.ID == 0 {
		err := cStore.Create(postEntity)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, &StandardResult{
			Message: "新增成功!",
			Data:    postEntity.ID,
		})
	}

	//edit
	err = cStore.Update(postEntity)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &StandardResult{
		Message: "修改成功!",
		Data:    postEntity.ID,
	})
}

func deleteColumn(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)
	id, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		return err
	}
	cStore := store.Stores().ColumnStore
	err = cStore.Delete(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &StandardResult{
		Message: "删除成功",
	})
}
