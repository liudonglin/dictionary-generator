package tpl

import (
	"dg-server/core"
	"dg-server/store"
	"strings"

	"github.com/flosch/pongo2"
)

var tableTemplete = `
using System;

namespace entity
{
	public class {{ fn.ToCamelString(table.Name) }} 
	{
		{% for column in columns %}
		public {{ fn.SqlTypeConvertLanguageType(column,project.DataBase,project.Language) }} {{ fn.ToCamelString(column.Name) }} { get; set; }
		{% endfor %}
	}
}
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

	out, err := tpl.Execute(pongo2.Context{
		"project":  project,
		"database": database,
		"table":    table,
		"columns":  columns,
		"indexs":   indexs,
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
func GetTableScript(req *core.TempleteLoadReq) (string, error) {
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

	out, err := tpl.Execute(pongo2.Context{
		"project":  project,
		"database": database,
		"table":    table,
		"columns":  columns,
		"indexs":   indexs,
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
