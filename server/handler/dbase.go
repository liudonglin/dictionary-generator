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

func saveDataBase(c echo.Context) error {
	postEntity := &core.DataBase{}
	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, postEntity)

	//验证字段
	errs := validator.ValidateStruct(postEntity)
	if errs != nil {
		return &BusinessError{Message: errs.Error()}
	}

	dbStore := store.Stores().DataBaseStore

	//检查名称是否重复
	dbEntity, err := dbStore.FindNameAndPID(postEntity.PID, postEntity.Name)
	if err != nil {
		return err
	}
	if dbEntity.ID != 0 && dbEntity.ID != postEntity.ID {
		return &BusinessError{Message: fmt.Sprintf("名称为%s的数据库已存在!", postEntity.Name)}
	}

	//add
	if postEntity.ID == 0 {
		if postEntity.PID == 0 {
			return &BusinessError{Message: "新增时必须传入项目id!"}
		}

		projectStore := store.Stores().ProjectStore
		dbProject, _ := projectStore.FindID(postEntity.PID)
		if dbProject.ID == 0 {
			return &BusinessError{Message: fmt.Sprintf("项目id: %d无效!", postEntity.PID)}
		}

		err := dbStore.Create(postEntity)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, &StandardResult{
			Message: "新增成功!",
			Data:    postEntity.ID,
		})
	}

	//edit
	err = dbStore.Update(postEntity)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &StandardResult{
		Message: "修改成功!",
		Data:    postEntity.ID,
	})
}

func listDataBase(c echo.Context) error {
	q := &core.DBQuery{}
	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, q)

	dbStore := store.Stores().DataBaseStore
	list, total, err := dbStore.List(q)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &StandardResult{
		Data: &PageResult{Total: total, List: list},
	})
}

func deleteDataBase(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)
	id, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		return err
	}

	//删除数据列
	cStore := store.Stores().ColumnStore
	err = cStore.DeleteByDID(id)
	if err != nil {
		return err
	}

	//删除数据表
	tableStore := store.Stores().TableStore
	err = tableStore.DeleteByDID(id)
	if err != nil {
		return err
	}

	//删除数据库
	dbStore := store.Stores().DataBaseStore
	err = dbStore.Delete(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &StandardResult{
		Message: "删除成功",
	})
}
