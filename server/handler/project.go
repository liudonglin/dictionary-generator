package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"code-server/core"
	"code-server/store"
	"code-server/validator"

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

	//检查项目名称是否重复
	dbProject, err := projectStore.FindName(project.Name)
	if err != nil {
		return err
	}
	if dbProject.ID != 0 && dbProject.ID != project.ID {
		return &BusinessError{Message: fmt.Sprintf("名称为%s的项目已存在!", project.Name)}
	}

	if project.ID == 0 {
		err := projectStore.Create(project)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, &StandardResult{
			Message: "新增成功!",
			Data:    project.ID,
		})
	}

	err = projectStore.Update(project)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &StandardResult{
		Message: "修改成功!",
		Data:    project.ID,
	})
}

func listProject(c echo.Context) error {
	project := &core.Project{}
	body, _ := ioutil.ReadAll(c.Request().Body)
	json.Unmarshal(body, project)

	projectStore := store.Stores().ProjectStore
	list, err := projectStore.List(project.Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, &StandardResult{
		Data: list,
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
