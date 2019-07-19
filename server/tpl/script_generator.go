package tpl

import (
	"dg-server/core"
	"dg-server/store"
	"strings"
	"time"

	"github.com/flosch/pongo2"
)

var tableTemplete = `
<?xml version="1.0" encoding="UTF-8" ?>
<!DOCTYPE mapper PUBLIC "-//mybatis.org//DTD Mapper 3.0//EN" "http://mybatis.org/dtd/mybatis-3-mapper.dtd" >
<mapper namespace="{{ project.NameSpace }}.dao.{{ fn.ToCamelString(table.Name) }}">
    
    <select id="load" resultType="{{ project.NameSpace }}.entity.{{ fn.ToCamelString(table.Name) }}">
        SELECT {% for column in columns %}
        {{column.Name}}{% if !fn.IsLastColumn(column,columns) %},{% endif %}{% endfor %}
        FROM {{ table.Name }}
        WHERE {{ fn.GetPK(indexs) }} = #{id}
    </select>

</mapper>
`

// TestGetTableScript ...
func TestGetTableScript(tid int64) string {
	tpl, err := pongo2.FromString(tableTemplete)

	tableStore := store.Stores().TableStore
	table, _ := tableStore.FindID(tid)

	dbStore := store.Stores().DataBaseStore
	database, _ := dbStore.FindID(table.DID)

	projectStore := store.Stores().ProjectStore
	project, _ := projectStore.FindID(table.PID)

	columnStore := store.Stores().ColumnStore
	columns, _, _ := columnStore.List(&core.ColumnQuery{
		TID: table.ID,
		Pager: core.Pager{
			Index: 0,
			Size:  9999999,
		},
	})

	indexs := make([]*core.Column, 0)
	for _, column := range columns {
		if column.Index {
			indexs = append(indexs, column)
		}
	}
	user := &core.User{
		Login: "liudonglin",
	}

	//解析枚举
	enums := getEnums(columns)

	out, err := tpl.Execute(pongo2.Context{
		"project":  project,
		"database": database,
		"table":    table,
		"columns":  columns,
		"indexs":   indexs,
		"user":     user,
		"enums":    enums,
		"fn":       fn,
	})

	if err != nil {
		print(err.Error())
	} else {
		print(out)
	}

	return out
}

// GetTableScript ...
func GetTableScript(req *core.TempleteLoadReq, userName string) (string, error) {
	tplStore := store.Stores().TempleteStore
	temp, _ := tplStore.FindID(req.TempleteID)
	tpl, err := pongo2.FromString(temp.Content)

	tableStore := store.Stores().TableStore
	table, _ := tableStore.FindID(req.TID)

	dbStore := store.Stores().DataBaseStore
	database, _ := dbStore.FindID(table.DID)

	projectStore := store.Stores().ProjectStore
	project, _ := projectStore.FindID(table.PID)

	columnStore := store.Stores().ColumnStore
	columns, _, _ := columnStore.List(&core.ColumnQuery{
		TID: table.ID,
		Pager: core.Pager{
			Index: 0,
			Size:  9999999,
		},
	})

	indexs := make([]*core.Column, 0)
	for _, column := range columns {
		if column.Index {
			indexs = append(indexs, column)
		}
	}

	user := &core.User{
		Login: userName,
	}

	//解析枚举
	enums := getEnums(columns)

	out, err := tpl.Execute(pongo2.Context{
		"project":  project,
		"database": database,
		"table":    table,
		"columns":  columns,
		"indexs":   indexs,
		"user":     user,
		"enums":    enums,
		"fn":       fn,
	})

	if err != nil {
		return "", err
	}
	return out, nil
}

// FnWrap ...
type FnWrap struct{}

var fn = &FnWrap{}

// LenColumn 获取列长度
func (*FnWrap) LenColumn(arr []*core.Column) int {
	return len(arr)
}

// IsLastColumn 判断是否是最后一列
func (*FnWrap) IsLastColumn(item *core.Column, list []*core.Column) bool {
	if item.ID == list[len(list)-1].ID {
		return true
	}
	return false
}

// IsLastKvd ...
func (*FnWrap) IsLastKvd(item *KVD, list []*KVD) bool {
	if item.Key == list[len(list)-1].Key {
		return true
	}
	return false
}

// ToSnakeString snake string, XxYy to xx_yy , XxYY to xx_yy
func (*FnWrap) ToSnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// ToCamelString camel string, xx_yy to XxYy
func (*FnWrap) ToCamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// FirstToLower 首字母转小写
func (*FnWrap) FirstToLower(s string) string {
	if len(s) < 1 {
		return ""
	}
	strArry := []rune(s)
	if strArry[0] >= 65 && strArry[0] <= 90 {
		strArry[0] += 32
	}
	return string(strArry)
}

// FirstToUpper 首字母转大写
func (*FnWrap) FirstToUpper(s string) string {
	if len(s) < 1 {
		return ""
	}
	strArry := []rune(s)
	if strArry[0] >= 97 && strArry[0] <= 122 {
		strArry[0] -= 32
	}
	return string(strArry)
}

// GetNowDate 获取当前日期
func (*FnWrap) GetNowDate() string {
	return time.Now().Format("2006-01-02")
}

// GetPK 获取主键
func (*FnWrap) GetPK(cols []*core.Column) *core.Column {
	for _, col := range cols {
		if col.PK {
			return col
		}
	}
	return nil
}

// ToUpper 转大写
func (*FnWrap) ToUpper(s string) string {
	return strings.ToUpper(s)
}

// ToLower 转小写
func (*FnWrap) ToLower(s string) string {
	return strings.ToLower(s)
}
