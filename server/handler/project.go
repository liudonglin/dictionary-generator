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

	//为项目匹配模版
	tq := &core.TempleteQuery{
		Pager: core.Pager{
			Index: 0,
			Size:  9999999,
		},
	}
	tplStore := store.Stores().TempleteStore
	tpls, _, err := tplStore.List(tq)
	if err != nil {
		return err
	}
	tplStore.DeleteProjectTempleteRelationByPID(project.ID)
	for _, tpl := range tpls {
		if matchTemplete(project, tpl) {
			tplStore.CreateProjectTempleteRelation(project.ID, tpl.ID)
		}
	}

	return c.JSON(http.StatusOK, &StandardResult{
		Message: msg,
		Data:    project.ID,
	})
}

func matchTemplete(p *core.Project, t *core.Templete) bool {
	matchDataBase, matchLanguage, matchOrm := false, false, false

	if t.DataBase == "" || t.DataBase == p.DataBase {
		matchDataBase = true
	}
	if t.Language == "" || t.Language == p.Language {
		matchLanguage = true
	}
	if t.Orm == "" || t.Orm == p.Orm {
		matchOrm = true
	}

	if matchDataBase && matchLanguage && matchOrm {
		return true
	}
	return false
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

	//删除数据列
	cStore := store.Stores().ColumnStore
	err = cStore.DeleteByPID(id)
	if err != nil {
		return err
	}

	//删除数据表
	tableStore := store.Stores().TableStore
	err = tableStore.DeleteByPID(id)
	if err != nil {
		return err
	}

	//删除数据库
	dbStore := store.Stores().DataBaseStore
	err = dbStore.DeleteByPID(id)
	if err != nil {
		return err
	}

	//删除连接信息
	connStore := store.Stores().ConnectionStore
	err = connStore.DeleteByPID(id)
	if err != nil {
		return err
	}

	//删除项目
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
