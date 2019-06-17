package handler

import (
	"code-server/core"
	"code-server/store"
	"code-server/validator"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

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

	postEntity.Created = time.Now().Format("2006-01-02 15:04:05")
	postEntity.Updated = postEntity.Created
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
