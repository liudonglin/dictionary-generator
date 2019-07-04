package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"dg-server/core"
	"dg-server/store"
	"dg-server/validator"

	"github.com/labstack/echo"
)

func saveProject(c echo.Context) error {
	project := &core.Project{}
	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, project)

	//验证字段
	errs := validator.ValidateStruct(project)
	if errs != nil {
		return &BusinessError{Message: errs.Error()}
	}

	project.Created = time.Now().Format("2006-01-02 15:04:05")
	project.Updated = project.Created
	projectStore := store.Stores().ProjectStore
	connStore := store.Stores().ConnectionStore

	//检查项目名称是否重复
	dbProject, err := projectStore.FindName(project.Name)
	if err != nil {
		return err
	}
	if dbProject.ID != 0 && dbProject.ID != project.ID {
		return &BusinessError{Message: fmt.Sprintf("名称为%s的项目已存在!", project.Name)}
	}

	msg := ""

	if project.ID == 0 {
		err := projectStore.Create(project)
		if err != nil {
			return err
		}
		msg = "新增成功!"
	} else {
		err = projectStore.Update(project)
		if err != nil {
			return err
		}
		msg = "修改成功!"
	}

	//保存链接信息
	connStore.DeleteByPID(project.ID)
	if project.Connections != nil {
		for _, item := range project.Connections {
			item.PID = project.ID
			item.DataBase = project.DataBase
			connStore.Create(item)
		}
	}

	return c.JSON(http.StatusOK, &StandardResult{
		Message: msg,
		Data:    project.ID,
	})
}

func loadProject(c echo.Context) error {
	project := &core.Project{}
	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, project)

	projectStore := store.Stores().ProjectStore

	dbProject, err := projectStore.FindID(project.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &StandardResult{
		Data: dbProject,
	})
}

func listProject(c echo.Context) error {
	q := &core.ProjectQuery{}
	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, q)

	projectStore := store.Stores().ProjectStore
	list, total, err := projectStore.List(q)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &StandardResult{
		Data: &PageResult{Total: total, List: list},
	})
}

func deleteProject(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)
	id, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		return err
	}
	projectStore := store.Stores().ProjectStore
	err = projectStore.Delete(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &StandardResult{
		Message: "删除成功",
	})
}

func loadConnectionsByPID(c echo.Context) error {
	body, _ := ioutil.ReadAll(c.Request().Body)
	id, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		return err
	}
	connStore := store.Stores().ConnectionStore

	q := &core.ConnectionQuery{
		PID: id,
		Pager: core.Pager{
			Index: 0,
			Size:  999999,
		},
	}

	list, _, err := connStore.List(q)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &StandardResult{
		Data: list,
	})
}
