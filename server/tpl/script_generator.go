package tpl

import (
	"dg-server/core"
	"dg-server/store"
	"strings"

	"github.com/flosch/pongo2"
)

var tableTemplete = `
package entity

type {{ toCamelString(table.Name) }} struct { {% for column in columns %}
	{{ toCamelString(column.Name) }} string ` + "`" + `json:"{{ toSnakeString(column.Name) }}"` + "`" + `{% endfor %}
}
`

// TestGetTableScript ...
func TestGetTableScript(tid int64) string {
	tpl, err := pongo2.FromString(tableTemplete)

	tableStore := store.Stores().TableStore
	table, err := tableStore.FindID(tid)
	if err != nil {
		return ""
	}

	dbStore := store.Stores().DataBaseStore
	database, err := dbStore.FindID(table.DID)
	if err != nil {
		return ""
	}

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
		"database":      database,
		"table":         table,
		"columns":       columns,
		"indexs":        indexs,
		"lenColumn":     lenColumn,
		"isLastColumn":  isLastColumn,
		"toCamelString": toCamelString,
		"toSnakeString": toSnakeString,
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
	table, err := tableStore.FindID(req.TID)
	if err != nil {
		return "", err
	}

	dbStore := store.Stores().DataBaseStore
	database, err := dbStore.FindID(table.DID)
	if err != nil {
		return "", err
	}

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
		"database":      database,
		"table":         table,
		"columns":       columns,
		"indexs":        indexs,
		"lenColumn":     lenColumn,
		"isLastColumn":  isLastColumn,
		"toCamelString": toCamelString,
		"toSnakeString": toSnakeString,
	})

	if err != nil {
		return "", err
	}

	return out, nil
}

func lenColumn(arr []*core.Column) int {
	return len(arr)
}

func isLastColumn(item *core.Column, list []*core.Column) bool {
	if item.ID == list[len(list)-1].ID {
		return true
	}
	return false
}

// snake string, XxYy to xx_yy , XxYY to xx_yy
func toSnakeString(s string) string {
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

// camel string, xx_yy to XxYy
func toCamelString(s string) string {
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
