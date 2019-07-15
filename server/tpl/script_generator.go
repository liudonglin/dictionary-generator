package tpl

import (
	"dg-server/core"
	"dg-server/store"

	"github.com/flosch/pongo2"
)

var tableTemplete = `
USE {{database.Name}} ;
Drop Table If Exists {{table.Name}} ;
CREATE TABLE {{table.Name}} ( {% for column in columns %}
{{column.Name}} {{column.ColumnType}} {% if column.Null %}{% else %}NOT{% endif %} NULL {% if column.AI %}AUTO_INCREMENT{% endif %} COMMENT '{{column.Title}}' {% if !isLastColumn(column,columns) %},{% endif %}{% endfor %}{% if lenColumn(indexs)>0 %},{% endif %}
{% for column in indexs %}
{% if column.PK %}PRIMARY {% elif column.Unique %}UNIQUE {% endif %}KEY {% if !column.PK %}idx_{{table.Name}}_{{column.Name}}{% endif %} ({{column.Name}}) {% if !isLastColumn(column,indexs) %},{% endif %} {% endfor %}
)ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='{{database.Description}}';
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
		"database":     database,
		"table":        table,
		"columns":      columns,
		"indexs":       indexs,
		"lenColumn":    lenColumn,
		"isLastColumn": isLastColumn,
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
		"database":     database,
		"table":        table,
		"columns":      columns,
		"indexs":       indexs,
		"lenColumn":    lenColumn,
		"isLastColumn": isLastColumn,
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
